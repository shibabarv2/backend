package routes

import (
	"github.com/gofiber/fiber/v2"
	"shiba-backend/structs"
)

func Index(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "OK",
		"errors":  structs.Errors{},
		"message": "Hi",
	})
}
