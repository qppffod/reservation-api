package api

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/qppffod/reservation-api/db"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		headers, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			fmt.Errorf("unauthorized")
		}

		if len(headers) < 1 {
			fmt.Println("Invalid headers, headers array length should be at least 1")
		}

		token := headers[0]

		claims, err := Parse(token)
		if err != nil {
			return err
		}

		expires := claims["expires"].(float64)
		if time.Now().Unix() > int64(expires) {
			return fmt.Errorf("token expired")
		}

		userID := claims["id"].(string)

		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			return err
		}

		c.Context().SetUserValue("user", user)

		return c.Next()
	}
}

func Parse(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method ", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("Failed to parse JWT Token:", err)
		return nil, fmt.Errorf("unauthorized")
	}

	if !token.Valid {
		return nil, fmt.Errorf("unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	return claims, nil
}
