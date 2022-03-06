package routes

import "github.com/gofiber/fiber/v2"

func NotFound(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"status":  "ERROR",
		"errors":  errors{"ROUTE_NOT_FOUND"},
		"message": "This route could not be found. Please provide a valid route and try again.",
	})
}
