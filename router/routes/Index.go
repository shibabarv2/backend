package routes

import (
	"github.com/gofiber/fiber/v2"
	"shiba-backend/structs"
)

func Index(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"status":  "OK",
		"errors":  structs.Errors{},
		"message": "Hi",
	})
}
