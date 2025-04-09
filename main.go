package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/qppffod/reservation-api/api"
	"github.com/qppffod/reservation-api/api/apiError"
	"github.com/qppffod/reservation-api/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: apiError.ErrorHandler,
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
		userStore    = db.NewMongoUserStore(client, db.DBNAME)
		hotelStore   = db.NewMongoHotelStore(client, db.DBNAME)
		roomStore    = db.NewMongoRoomStore(client, db.DBNAME, hotelStore)
		bookingStore = db.NewMongoBookingStore(client, db.DBNAME)
		store        = &db.Store{
			User:    userStore,
			Hotel:   hotelStore,
			Room:    roomStore,
			Booking: bookingStore,
		}

		// handlers
		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		authHandler    = api.NewAuthHandler(userStore)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)

		app   = fiber.New(config)
		auth  = app.Group("/api")
		apiv1 = app.Group("/api/v1", api.JWTAuthentication(userStore))
		admin = apiv1.Group("/admin", api.AdminAuth)
	)

	// auth
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// user handlers
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)

	// hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRoomsByHotelID)

	// room handlers
	apiv1.Get("/room", roomHandler.HandleGetRooms)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	// booking handlers
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	// admin
	admin.Get("/booking", bookingHandler.HandleGetBookings)

	log.Printf("Listening on port %s\n", *listenAddr)
	app.Listen(*listenAddr)
}
