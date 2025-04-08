package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/qppffod/reservation-api/api"
	"github.com/qppffod/reservation-api/db"
	"github.com/qppffod/reservation-api/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	hotelStore   db.HotelStore
	roomStore    db.RoomStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	ctx          = context.Background()
)

func seedUser(isAdmin bool, fname, lname, email, password string) *types.User {
	user, err := types.NewUserFromParams(&types.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  password,
	})
	user.IsAdmin = isAdmin
	if err != nil {
		log.Fatal(err)
	}

	insertedUser, err := userStore.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))
	return insertedUser
}

func seedRoom(size string, ss bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: ss,
		Price:   price,
		HotelID: hotelID,
	}

	insertedRoom, err := roomStore.InsertRoom(ctx, room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func seedBooking(userID, roomID primitive.ObjectID, numPers int, from, till time.Time) *types.Booking {
	booking := &types.Booking{
		UserID:     userID,
		RoomID:     roomID,
		NumPersons: numPers,
		FromDate:   from,
		TillDate:   till,
	}

	insertedBooking, err := bookingStore.InsertBooking(ctx, booking)
	if err != nil {
		log.Fatal(err)
	}

	return insertedBooking
}

func seedHotel(name, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func main() {
	user := seedUser(false, "John", "Smith", "john@smith.com", "userpassword")
	admin := seedUser(true, "admin", "admin", "admin@admin.com", "adminpassword")
	_ = admin

	seedHotel("Bellucia", "France", 3)
	seedHotel("Legoland", "Denmark", 4)
	hotel := seedHotel("Random Hotel", "Spain", 1)

	seedRoom("small", true, 89.99, hotel.ID)
	seedRoom("medium", true, 189.99, hotel.ID)
	room := seedRoom("large", true, 289.99, hotel.ID)

	booking := seedBooking(user.ID, room.ID, 5, time.Now(), time.Now().AddDate(0, 0, 2))
	fmt.Printf("booking -> %s\n", booking.ID.Hex())
}

func init() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client, db.DBNAME)
	roomStore = db.NewMongoRoomStore(client, db.DBNAME, hotelStore)
	userStore = db.NewMongoUserStore(client, db.DBNAME)
	bookingStore = db.NewMongoBookingStore(client, db.DBNAME)
}
