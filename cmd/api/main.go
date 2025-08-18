package main

import (
	"Go-Microservice/internal/env"
	"log"
)

func main() {

	app := &application{
		config: config{
			addr: env.GetPort("ADDR", ":8080"),
		},
	}

	mount := app.mount()
	log.Fatal(app.run(mount))
}
