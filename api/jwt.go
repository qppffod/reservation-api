package api

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	headers, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		fmt.Errorf("unauthorized")
	}

	if len(headers) < 1 {
		fmt.Println("Invalid headers, headers array length should be at least 1")
	}

	token := headers[0]

	err := Parse(token)
	if err != nil {
		return err
	}

	fmt.Println(token)

	return nil
}

func Parse(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Printf("invalid signing method: %s", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("Failed to parse JWT Token:", err)
		return fmt.Errorf("unauthorized")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims)
	}

	return fmt.Errorf("unauthorized")
}
