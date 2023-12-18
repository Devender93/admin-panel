package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kominkamen/rootds-admin/models"
)

var (
	accessTokenSecret = []byte("KXsMPri4PLlFlGcqU0f4P9y2s0aIOos9")
)

func ValidateAuthToken(c *fiber.Ctx) error {
	tokenString := c.Get("X-Access-Token")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Error Token is Empty",
			StatusCode: fiber.StatusUnauthorized,
		})
	}

	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return accessTokenSecret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Error Invalid Token",
			StatusCode: fiber.StatusUnauthorized,
		})
	}

	customClaims, ok := token.Claims.(*models.Claims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Error Invalid Claims",
			StatusCode: fiber.StatusUnauthorized,
		})
	}

	if customClaims.Role == nil || *customClaims.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(models.JsonResponse{
			Status:     false,
			Message:    "Error: You are Unauthorized",
			StatusCode: fiber.StatusForbidden,
		})
	}
	return c.Next()
}
