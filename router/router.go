package router

import (
	"github.com/gofiber/fiber/v2"
	"shiba-backend/router/routes"
)

func Router(server *fiber.App) {
	group := server.Group("/v2")
	group.Get("/", routes.Index)
	group.Get("/stats", routes.Stats)

	auth := group.Group("/auth")
	auth.Get("/register", routes.Register)

	// * Keep this at the end
	server.Get("*", routes.NotFound)
	server.Post("*", routes.NotFound)
}
