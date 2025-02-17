package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"

	"github.com/ariefro/threads-server/docs"
	"github.com/ariefro/threads-server/internal/auth"
	"github.com/ariefro/threads-server/internal/mailer"
	"github.com/ariefro/threads-server/internal/ratelimiter"
	"github.com/ariefro/threads-server/internal/store"
	"github.com/ariefro/threads-server/internal/store/cache"
)

type application struct {
	config        config
	store         store.Storage
	cacheStorage  cache.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.EmailSender
	authenticator auth.Authenticator
	rateLimiter   ratelimiter.Limiter
}

type config struct {
	addr              string
	db                dbConfig
	env               string
	mail              mailConfig
	frontendURL       string
	auth              authConfig
	redisCfg          redisConfig
	rateLimiter       ratelimiter.Config
	corsAllowedOrigin string
	sentry            sentryConfig
}

type redisConfig struct {
	addr     string
	password string
	db       int
	enabled  bool
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type basicConfig struct {
	user string
	pass string
}

type mailConfig struct {
	fromName  string
	fromEmail string
	password  string
	exp       time.Duration
}

type dbConfig struct {
	driver       string
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type sentryConfig struct {
	dsn         string
	sampleRate  float64
	environment string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{app.config.corsAllowedOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	if app.config.rateLimiter.Enabled {
		r.Use(app.RateLimiterMiddleware)
	}

	sentryHandler := sentryhttp.New(sentryhttp.Options{
		Repanic:         true, // Allow panic to be propagated after Sentry captures it
		WaitForDelivery: true, // Block the request until Sentry confirms the event was sent
		Timeout:         2,    // Set timeout for event delivery in seconds
	})

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		r.With(app.BasicAuthMiddleware()).Get("/debug/vars", expvar.Handler().ServeHTTP)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/", sentryHandler.HandleFunc(app.createPostHandler))

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)

				r.Get("/", sentryHandler.HandleFunc(app.getPostHandler))
				r.Patch("/", sentryHandler.HandleFunc(app.checkPostOwnership("moderator", app.updatePostHandler)))
				r.Delete("/", sentryHandler.HandleFunc(app.checkPostOwnership("admin", app.deletePostHandler)))

				r.Post("/comment", sentryHandler.HandleFunc(app.createCommentHandler))
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", sentryHandler.HandleFunc(app.activateUserHandler))

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)

				r.Get("/", sentryHandler.HandleFunc(app.getUserHandler))
				r.Put("/follow", sentryHandler.HandleFunc(app.followUserHandler))
				r.Put("/unfollow", sentryHandler.HandleFunc(app.unfollowUserHandler))
			})

			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/feed", sentryHandler.HandleFunc(app.getUserFeedHandler))
			})
		})

		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", sentryHandler.HandleFunc(app.registerUserHandler))
			r.Post("/token", sentryHandler.HandleFunc(app.createTokenHandler))
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	// Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.frontendURL
	docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Infow("signal caught", "signal", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	app.logger.Infow("server has started", "addr", app.config.addr, "env", app.config.env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		sentry.CaptureException(err)
		return err
	}

	err = <-shutdown
	if err != nil {
		sentry.CaptureException(err)
		return err
	}

	app.logger.Infow("server has stopped", "addr", app.config.addr, "env", app.config.env)

	return nil
}
