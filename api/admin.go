package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/qppffod/reservation-api/types"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return fmt.Errorf("unauthorized")
	}

	if !user.IsAdmin {
		return fmt.Errorf("error: access denied. only for admins")
	}

	return c.Next()
}
