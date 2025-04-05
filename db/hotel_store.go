package db

import (
	"context"

	"github.com/qppffod/reservation-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelStore interface {
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	UpdateHotel(context.Context, *types.Room) error
	GetHotels(context.Context) ([]*types.Hotel, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client, dbname string) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(dbname).Collection("hotels"),
	}
}

func (s *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}

	hotel.ID = res.InsertedID.(primitive.ObjectID)

	return hotel, nil
}

func (s *MongoHotelStore) UpdateHotel(ctx context.Context, room *types.Room) error {
	update := bson.M{"$push": bson.M{"rooms": room.ID}}

	_, err := s.coll.UpdateByID(ctx, room.HotelID, update)
	return err
}

func (s *MongoHotelStore) GetHotels(ctx context.Context) ([]*types.Hotel, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var hotels []*types.Hotel
	if err := cur.All(ctx, &hotels); err != nil {
		return nil, err
	}

	return hotels, nil
}
