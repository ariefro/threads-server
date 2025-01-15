package main

import (
	"fmt"
	"log"

	"github.com/ariefro/threads-server/internal/db"
	"github.com/ariefro/threads-server/internal/env"
	"github.com/ariefro/threads-server/internal/repository"
)

func main() {
	env, err := env.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	cfg := config{
		addr: env.Server.Port,
		db: dbConfig{
			driver: env.Database.Driver,
			dsn: fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
				env.Database.Host,
				env.Database.Port,
				env.Database.User,
				env.Database.Password,
				env.Database.DBName,
				env.Database.SSLMode,
			),
			maxOpenConns: env.Database.MaxOpenConn,
			maxIdleConns: env.Database.MaxIdleConn,
			maxIdleTime:  env.Database.MaxIdleTime,
		},
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

	repository := repository.NewRepositories(db)

	app := &application{
		config:     cfg,
		repository: *repository,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
