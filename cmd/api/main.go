package main

import (
	"log"

	"github.com/ariefro/threads-server/internal/env"
)

func main() {
	env, err := env.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	cfg := config{
		addr: env.Server.Port,
	}

	app := &application{
		config: cfg,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
