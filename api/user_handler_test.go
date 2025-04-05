package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/qppffod/reservation-api/db"
	"github.com/qppffod/reservation-api/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	return &testdb{
		UserStore: db.NewMongoUserStore(client, db.TESTDBNAME),
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	u := types.CreateUserParams{
		FirstName: "TestFirstName",
		LastName:  "TestLastName",
		Email:     "test@test.com",
		Password:  "12312311",
	}
	b, err := json.Marshal(u)
	if err != nil {
		t.Fatal("error marshal user")
	}

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	var user types.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Fatal(err)
	}

	if len(user.ID) == 0 {
		t.Fatalf("expected the user ID to be set")
	}
	if len(user.EncryptedPassword) > 0 {
		t.Fatalf("expected the EncryptedPassword not to be included in the JSON response")
	}
	if user.FirstName != u.FirstName {
		t.Fatalf("expected firstName to be %s but got %s", u.FirstName, user.FirstName)
	}
	if user.LastName != u.LastName {
		t.Fatalf("expected lastName to be %s but got %s", u.LastName, user.LastName)
	}
	if user.Email != u.Email {
		t.Fatalf("expected email to be %s but got %s", u.Email, user.Email)
	}

	fmt.Println("USER--->", user)
}
