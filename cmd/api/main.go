package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ariefro/threads-server/internal/auth"
	"github.com/ariefro/threads-server/internal/db"
	"github.com/ariefro/threads-server/internal/env"
	"github.com/ariefro/threads-server/internal/mailer"
	"github.com/ariefro/threads-server/internal/ratelimiter"
	"github.com/ariefro/threads-server/internal/store"
	"github.com/ariefro/threads-server/internal/store/cache"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const version = "0.0.1"

func main() {
	env, err := env.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	cfg := config{
		addr: env.AppPort,
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
		env: env.AppEnv,
		mail: mailConfig{
			exp:       time.Hour * 24 * 3, // 3 days
			fromName:  env.SenderName,
			fromEmail: env.EmailSender,
			password:  env.EmailSenderPassword,
		},
		frontendURL: env.FrontendURL,
		auth: authConfig{
			basic: basicConfig{
				user: env.AuthBasicUser,
				pass: env.AuthBasicPass,
			},
		},
		rateLimiter: ratelimiter.Config{
			RequestsPerTimeFrame: env.RateLimiterRequestCount,
			TimeFrame:            time.Second * 5,
			Enabled:              env.RateLimiterEnabled,
		},
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database
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

	store := store.NewStorage(db)
	cacheStorage := cache.NewRedisStorage(rdb)

	mailer := mailer.NewGmailSender(cfg.mail.fromName, cfg.mail.fromEmail, cfg.mail.password)

	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.token.secret,
		cfg.auth.token.iss,
		cfg.auth.token.iss,
	)

	app := &application{
		config:        cfg,
		store:         *store,
		cacheStorage:  *cacheStorage,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuthenticator,
		rateLimiter:   rateLimiter,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
