package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/qppffod/reservation-api/api/apiError"
	"github.com/qppffod/reservation-api/types"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return apiError.ErrUnAuthorized()
	}

	if !user.IsAdmin {
		return apiError.NewError(http.StatusUnauthorized, "error: access denied. only for admins")
	}

	return c.Next()
}
