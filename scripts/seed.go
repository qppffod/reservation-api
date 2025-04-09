package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/qppffod/reservation-api/api"
	"github.com/qppffod/reservation-api/db"
	"github.com/qppffod/reservation-api/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client, db.DBNAME)
	store := db.Store{
		Hotel:   hotelStore,
		Room:    db.NewMongoRoomStore(client, db.DBNAME, hotelStore),
		User:    db.NewMongoUserStore(client, db.DBNAME),
		Booking: db.NewMongoBookingStore(client, db.DBNAME),
	}

	user := fixtures.AddUser(&store, false, "john", "smith")
	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(&store, true, "admin", "admin")
	fmt.Printf("%s -> %s\n", admin.Email, api.CreateTokenFromUser(admin))
	hotel := fixtures.AddHotel(&store, "Bellucia", "France", 3, nil)
	room := fixtures.AddRoom(&store, "small", true, 89.99, hotel.ID)
	booking := fixtures.AddBooking(&store, user.ID, room.ID, 5, time.Now(), time.Now().AddDate(0, 0, 2))
	fmt.Println("booking ->", booking.ID.Hex())
}
