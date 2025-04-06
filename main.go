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

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "listen address of http server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatal(err)
	}

	var (
		// stores
		userStore  = db.NewMongoUserStore(client, db.DBNAME)
		hotelStore = db.NewMongoHotelStore(client, db.DBNAME)
		roomStore  = db.NewMongoRoomStore(client, db.DBNAME, hotelStore)
		store      = &db.Store{
			User:  userStore,
			Hotel: hotelStore,
			Room:  roomStore,
		}

		// handlers
		userHandler  = api.NewUserHandler(userStore)
		hotelHandler = api.NewHotelHandler(store)
		authHandler  = api.NewAuthHandler(userStore)

		app   = fiber.New(config)
		auth  = app.Group("/api")
		apiv1 = app.Group("/api/v1", api.JWTAuthentication)
	)

	// auth
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// user routes
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)

	// hotel routes
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRoomsByHotelID)

	log.Printf("Listening on port %s\n", *listenAddr)
	app.Listen(*listenAddr)
}
