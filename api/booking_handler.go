package api

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/qppffod/reservation-api/db"
	"github.com/qppffod/reservation-api/types"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

// only for admins
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), nil)
	if err != nil {
		return err
	}

	return c.JSON(bookings)
}

// only users
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")

	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}

	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return fmt.Errorf("unauthorized")
	}

	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(GenericResponse{
			Type: "error",
			Msg:  "unauthorized",
		})
	}

	return c.JSON(booking)
}
