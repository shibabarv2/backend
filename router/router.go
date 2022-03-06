package router

import (
	"github.com/gofiber/fiber/v2"
	"shiba-backend/router/routes"
)

func Router(server *fiber.App) {
	group := server.Group("/v2")
	group.Get("/", routes.Index)

	// * Keep this at the end
	server.Get("*", routes.NotFound)
	server.Post("*", routes.NotFound)
}
