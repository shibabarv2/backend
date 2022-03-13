package routes

import (
	"github.com/gofiber/fiber/v2"
	"shiba-backend/structs"
	"strings"
)

func NotFound(ctx *fiber.Ctx) error {

	if strings.Contains(ctx.Path(), "v2") {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"ROUTE_NOT_FOUND"},
			"message": "This route could not be found. Please provide a valid route and try again.",
		})
	}

	return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"status":  "ERROR",
		"errors":  structs.Errors{"ROUTE_NOT_FOUND"},
		"message": "This route could not be found. Please provide a valid route and try again. (Maybe you meant /v2" + ctx.Path() + "?)",
	})
}
