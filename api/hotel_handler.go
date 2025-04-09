package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/qppffod/reservation-api/api/apiError"
	"github.com/qppffod/reservation-api/db"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")

	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), id)
	if err != nil {
		return apiError.ErrResourseNotFound("hotel")
	}

	return c.JSON(hotel)
}

func (h *HotelHandler) HandleGetRoomsByHotelID(c *fiber.Ctx) error {
	id := c.Params("id")

	rooms, err := h.store.Room.GetRoomsByHotelID(c.Context(), id)
	if err != nil {
		return apiError.ErrResourseNotFound("hotel rooms")
	}

	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	hotels, err := h.store.Hotel.GetHotels(c.Context())
	if err != nil {
		return apiError.ErrResourseNotFound("hotels")
	}

	return c.JSON(hotels)
}
