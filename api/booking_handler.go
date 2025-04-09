package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/qppffod/reservation-api/api/apiError"
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
		return apiError.ErrResourseNotFound("cancel booking")
	}

	user, err := getAuthUser(c)
	if err != nil {
		return apiError.ErrUnAuthorized()
	}

	if booking.UserID != user.ID {
		return apiError.ErrUnAuthorized()
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
		return apiError.ErrResourseNotFound("bookings")
	}

	return c.JSON(bookings)
}

// only users
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")

	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return apiError.ErrResourseNotFound("booking")
	}

	user, err := getAuthUser(c)
	if err != nil {
		return err
	}

	if booking.UserID != user.ID {
		return apiError.ErrUnAuthorized()
	}

	return c.JSON(booking)
}
