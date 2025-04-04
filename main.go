package main

import (
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/qppffod/reservation-api/api"
)

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "listen address of http server")
	flag.Parse()

	app := fiber.New()
	apiv1 := app.Group("/api/v1")

	apiv1.Get("/", api.HandleGetUser)

	log.Printf("Listening on port %s\n", *listenAddr)
	app.Listen(*listenAddr)
}
