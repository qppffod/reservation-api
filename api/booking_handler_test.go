package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/qppffod/reservation-api/api/apiError"
	"github.com/qppffod/reservation-api/db/fixtures"
	"github.com/qppffod/reservation-api/types"
)

func TestUserGetBooking(t *testing.T) {
	tdb := setup(t)
	tdb.teardown(t)

	var (
		nonAuthUser    = fixtures.AddUser(tdb.store, false, "Some", "Noname")
		user           = fixtures.AddUser(tdb.store, false, "john", "smith")
		hotel          = fixtures.AddHotel(tdb.store, "Bellucia", "France", 4, nil)
		room           = fixtures.AddRoom(tdb.store, "small", true, 89.99, hotel.ID)
		from           = time.Now()
		till           = time.Now().AddDate(0, 0, 2)
		booking        = fixtures.AddBooking(tdb.store, user.ID, room.ID, 2, from, till)
		app            = fiber.New()
		route          = app.Group("/:id", JWTAuthentication(tdb.store.User))
		bookingHandler = NewBookingHandler(tdb.store)
	)
	_ = booking

	route.Get("/", bookingHandler.HandleGetBooking)

	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response, got %d", resp.StatusCode)
	}

	var bookingResp types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}

	if bookingResp.ID != booking.ID {
		t.Fatalf("expected %s but got %s", booking.ID, bookingResp.ID)
	}
	if bookingResp.UserID != booking.UserID {
		t.Fatalf("expected %s but got %s", booking.UserID, bookingResp.UserID)
	}

	// test not auth user access
	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthUser))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected a non 200 response but got %d", resp.StatusCode)
	}
}

func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	tdb.teardown(t)

	var (
		config         = fiber.Config{ErrorHandler: apiError.ErrorHandler}
		adminUser      = fixtures.AddUser(tdb.store, true, "admin", "admin")
		user           = fixtures.AddUser(tdb.store, false, "john", "smith")
		hotel          = fixtures.AddHotel(tdb.store, "Bellucia", "France", 4, nil)
		room           = fixtures.AddRoom(tdb.store, "small", true, 89.99, hotel.ID)
		from           = time.Now()
		till           = time.Now().AddDate(0, 0, 2)
		booking        = fixtures.AddBooking(tdb.store, user.ID, room.ID, 2, from, till)
		app            = fiber.New(config)
		admin          = app.Group("/", JWTAuthentication(tdb.store.User), AdminAuth)
		bookingHandler = NewBookingHandler(tdb.store)
	)
	_ = booking
	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 response, got %d", resp.StatusCode)
	}

	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}

	if len(bookings) != 1 {
		t.Fatalf("expecteed 1 booking but got %d", len(bookings))
	}
	have := bookings[0]
	if have.ID != booking.ID {
		t.Fatalf("expected %s but got %s", booking.ID, have.ID)
	}
	if have.UserID != booking.UserID {
		t.Fatalf("expected %s but got %s", booking.UserID, have.UserID)
	}

	// test non-admin access
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected a non 200 response but got %d", resp.StatusCode)
	}

}
