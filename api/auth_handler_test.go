package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/qppffod/reservation-api/db"
	"github.com/qppffod/reservation-api/types"
)

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	tdb.teardown(t)

	insertedUser := seedTestUser(t, tdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	authParams := types.AuthParams{
		Email:    "test@test.com",
		Password: "testpassword",
	}
	b, err := json.Marshal(authParams)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http status of 200 but got %d", resp.StatusCode)
	}

	var authResp types.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}

	if authResp.Token == "" {
		t.Fatal("expected the JWT Token to be present in the auth response")
	}
	insertedUser.EncryptedPassword = ""

	if !reflect.DeepEqual(insertedUser, authResp.User) {
		t.Fatal("expected the user to be the inserted user")
	}
}

func TestAuthenticateWithWrongPassword(t *testing.T) {
	tdb := setup(t)
	tdb.teardown(t)

	seedTestUser(t, tdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	authParams := types.AuthParams{
		Email:    "test@test.com",
		Password: "testpasswordWRONG",
	}
	b, err := json.Marshal(authParams)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected http status of 400 but got %d", resp.StatusCode)
	}

	var genResp GenericResponse
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}

	if genResp.Type != "error" {
		t.Fatalf("expected gen response type to be error but got %s", genResp.Type)
	}
	if genResp.Msg != "invalid credentials" {
		t.Fatalf("expected gen response msg to be <invalid crednetials> but got %s", genResp.Msg)
	}
}

func seedTestUser(t *testing.T, userStore db.UserStore) *types.User {
	testUser, err := types.NewUserFromParams(&types.CreateUserParams{
		FirstName: "ftest",
		LastName:  "ltest",
		Email:     "test@test.com",
		Password:  "testpassword",
	})
	if err != nil {
		log.Fatal(err)
	}
	insertedUser, err := userStore.InsertUser(context.TODO(), testUser)
	if err != nil {
		log.Fatal(err)
	}

	return insertedUser
}
