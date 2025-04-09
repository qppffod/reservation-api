package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/qppffod/reservation-api/db"
	"github.com/qppffod/reservation-api/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(s *db.Store, admin bool, fname, lname string) *types.User {
	user, err := types.NewUserFromParams(&types.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     fmt.Sprintf("%s@%s.com", fname, lname),
		Password:  fmt.Sprintf("%s_%s", fname, lname),
	})
	user.IsAdmin = admin
	if err != nil {
		log.Fatal(err)
	}

	insertedUser, err := s.User.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	return insertedUser
}

func AddHotel(s *db.Store, name, location string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	var roomIDS = rooms
	if rooms == nil {
		roomIDS = []primitive.ObjectID{}
	}
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    roomIDS,
		Rating:   rating,
	}

	insertedHotel, err := s.Hotel.InsertHotel(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func AddRoom(s *db.Store, size string, ss bool, price float64, hID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: ss,
		Price:   price,
		HotelID: hID,
	}

	insertedRoom, err := s.Room.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func AddBooking(s *db.Store, uID, rID primitive.ObjectID, numPers int, from, till time.Time) *types.Booking {
	booking := &types.Booking{
		UserID:     uID,
		RoomID:     rID,
		NumPersons: numPers,
		FromDate:   from,
		TillDate:   till,
	}

	insertedBooking, err := s.Booking.InsertBooking(context.TODO(), booking)
	if err != nil {
		log.Fatal(err)
	}

	return insertedBooking
}
