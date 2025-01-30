package main

import (
	"expvar"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/ariefro/threads-server/internal/auth"
	"github.com/ariefro/threads-server/internal/db"
	"github.com/ariefro/threads-server/internal/env"
	"github.com/ariefro/threads-server/internal/mailer"
	"github.com/ariefro/threads-server/internal/ratelimiter"
	"github.com/ariefro/threads-server/internal/store"
	"github.com/ariefro/threads-server/internal/store/cache"
	"github.com/getsentry/sentry-go"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			Threads API
//	@description	API for Threads, a social network
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	env, err := env.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	cfg := config{
		addr: env.AppPort,
		env:  env.AppEnv,
		db: dbConfig{
			driver: env.DBDriver,
			dsn: fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
				env.DBHost,
				env.DBPort,
				env.DBUser,
				env.DBPassword,
				env.DBName,
				env.DBSSLMode,
			),
			maxOpenConns: env.DBMaxOpenConns,
			maxIdleConns: env.DBMaxIdleConns,
			maxIdleTime:  env.DBMaxIdleTime,
		},
		redisCfg: redisConfig{
			addr:     env.RedisAddress,
			password: env.RedisPassword,
			db:       env.RedisDB,
			enabled:  env.RedisEnabled,
		},
		auth: authConfig{
			basic: basicConfig{
				user: env.AuthBasicUser,
				pass: env.AuthBasicPass,
			},
			token: tokenConfig{
				secret: env.AuthTokenSecret,
				exp:    time.Hour * 24 * time.Duration(env.AuthTokenExpired),
			},
		},
		rateLimiter: ratelimiter.Config{
			RequestsPerTimeFrame: env.RateLimiterRequestCount,
			TimeFrame:            time.Second * 5,
			Enabled:              env.RateLimiterEnabled,
		},
		mail: mailConfig{
			exp:       time.Hour * 24 * 3, // 3 days
			fromName:  env.SenderName,
			fromEmail: env.EmailSender,
			password:  env.EmailSenderPassword,
		},
		sentry: sentryConfig{
			dsn:         env.SentryDsn,
			sampleRate:  env.SentrySampleRate,
			environment: env.AppEnv,
		},
		frontendURL:       env.FrontendURL,
		corsAllowedOrigin: env.CorsAllowedOrigin,
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Sentry
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.sentry.dsn,
		Environment:      cfg.sentry.environment,
		EnableTracing:    true,
		Debug:            true,
		TracesSampleRate: cfg.sentry.sampleRate,
	}); err != nil {
		logger.Errorw("sentry initialization failed", "error", err)
	}

	// Main Database
	db, err := db.NewDBConn(
		cfg.db.driver,
		cfg.db.dsn,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("database connection pool established")

	// cache
	var rdb *redis.Client
	if cfg.redisCfg.enabled {
		rdb = cache.NewRedisClient(cfg.redisCfg.addr, cfg.redisCfg.password, cfg.redisCfg.db)
		logger.Info("redis cache connection established")

		defer rdb.Close()
	}

	// rate limiter
	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.rateLimiter.RequestsPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	// mailer
	mailer := mailer.NewGmailSender(cfg.mail.fromName, cfg.mail.fromEmail, cfg.mail.password)

	// authenticator
	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.token.secret,
		cfg.auth.token.iss,
		cfg.auth.token.iss,
	)

	store := store.NewStorage(db)
	cacheStorage := cache.NewRedisStorage(rdb)

	app := &application{
		config:        cfg,
		store:         *store,
		cacheStorage:  *cacheStorage,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuthenticator,
		rateLimiter:   rateLimiter,
	}

	// metrics collected
	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
