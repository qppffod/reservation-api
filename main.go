package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/qppffod/reservation-api/api"
	"github.com/qppffod/reservation-api/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dburi = "mongodb://localhost:27017"

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "listen address of http server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatal(err)
	}

	userStore := db.NewMongoUserStore(client, db.DBNAME)

	userHandler := api.NewUserHandler(userStore)

	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")

	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)

	log.Printf("Listening on port %s\n", *listenAddr)
	app.Listen(*listenAddr)
}
