package main

import (
	"fmt"
	"log"

	"github.com/ariefro/threads-server/internal/db"
	"github.com/ariefro/threads-server/internal/env"
	"github.com/ariefro/threads-server/internal/store"
)

func main() {
	env, err := env.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		env.DBHost,
		env.DBPort,
		env.DBUser,
		env.DBPassword,
		env.DBName,
		env.DBSSLMode,
	)

	dbConn, err := db.NewDBConn(
		env.DBDriver,
		dsn,
		env.DBMaxOpenConns,
		env.DBMaxIdleConns,
		env.DBMaxIdleTime,
	)
	if err != nil {
		log.Fatal(err)
	}

	defer dbConn.Close()

	store := store.NewStorage(dbConn)

	db.Seed(*store, dbConn)
}
