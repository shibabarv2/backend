package routes

import "github.com/gofiber/fiber/v2"

type errors []string

func Index(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"status":  "OK",
		"errors":  errors{},
		"message": "Hi",
	})
}
