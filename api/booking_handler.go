package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/qppffod/reservation-api/db"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}

	user, err := getAuthUser(c)
	if err != nil {
		return err
	}

	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(GenericResponse{
			Type: "error",
			Msg:  "unauthorized",
		})
	}

	if err := h.store.Booking.UpdateBooking(c.Context(), id); err != nil {
		return err
	}

	return c.JSON(GenericResponse{
		Type: "msg",
		Msg:  "updated",
	})
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

	user, err := getAuthUser(c)
	if err != nil {
		return err
	}

	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(GenericResponse{
			Type: "error",
			Msg:  "unauthorized",
		})
	}

	return c.JSON(booking)
}
