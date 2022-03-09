package router

import (
	"github.com/gofiber/fiber/v2"
	"shiba-backend/router/routes"
	routes2 "shiba-backend/router/routes/admin"
	routes3 "shiba-backend/router/routes/auth"
	routes4 "shiba-backend/router/routes/user"
	//. "github.com/bestmethod/go-logger"
)

//LoggerP := gologg.new(gologg.Logger)

func Router(server *fiber.App) {

	group := server.Group("/v2")
	group.Get("/", routes.Index)
	group.Get("/stats", routes.Stats)

	auth := group.Group("/auth")
	auth.Get("/register", routes3.Register)

	admin := group.Group("/admin")
	admin.Get("/create/invite", routes2.AddInvite)
	admin.Get("/remove/invite", routes2.RemoveInvite)
	admin.Get("/blacklist/user", routes2.BanUser)
	admin.Get("/lookup/invite", routes2.GetInvite)

	user := group.Group("/users")
	user.Get("/lookup/:username", routes4.Lookup)

	// * Keep this at the end
	server.Get("*", routes.NotFound)
	server.Post("*", routes.NotFound)
}
