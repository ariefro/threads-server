package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ariefro/threads-server/internal/auth"
	"github.com/ariefro/threads-server/internal/db"
	"github.com/ariefro/threads-server/internal/env"
	"github.com/ariefro/threads-server/internal/mailer"
	"github.com/ariefro/threads-server/internal/store"
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

	store := store.NewStorage(db)

	mailer := mailer.NewGmailSender(cfg.mail.fromName, cfg.mail.fromEmail, cfg.mail.password)

	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.token.secret,
		cfg.auth.token.iss,
		cfg.auth.token.iss,
	)

	app := &application{
		config:        cfg,
		store:         *store,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuthenticator,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
