package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"shiba-backend/router"
	"shiba-backend/structs"
	"time"
)

func main() {
	godotenv.Load()
	key := os.Getenv("MONGODBURI")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(key))

	if err != nil {
		println("There was an error connecting to mongoDB")
		return
	}

	structs.DB = client.Database("shiba")

	// Server
	server := fiber.New()

	server.Use(logger.New())
	server.Use(limiter.New(limiter.Config{
		Max:        5,
		Expiration: 40 * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"status":      "ERROR",
				"errors":      structs.Errors{"SENDING_TOO_MANY_REQUESTS"},
				"message":     "Slow down! You're sending too many requests in such a short period of time. Please try again later.",
				"retry-after": c.GetRespHeader("Retry-After") + "s",
			})
		},
	}))

	router.Router(server)

	server.Listen(":3000")
}
