package main

import (
	"fmt"
	"log"

	"github.com/ariefro/threads-server/internal/db"
	"github.com/ariefro/threads-server/internal/env"
	"github.com/ariefro/threads-server/internal/store"
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
	}

	db, err := db.NewDBConn(
		cfg.db.driver,
		cfg.db.dsn,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Println("database connection pool established")

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  *store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
