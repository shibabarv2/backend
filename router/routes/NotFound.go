package routes

import (
	"github.com/gofiber/fiber/v2"
	"shiba-backend/structs"
)

func NotFound(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"status":  "ERROR",
		"errors":  structs.Errors{"ROUTE_NOT_FOUND"},
		"message": "This route could not be found. Please provide a valid route and try again.",
	})
}
