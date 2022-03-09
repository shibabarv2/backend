package routes

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"shiba-backend/structs"
	"shiba-backend/util"
)

func BanUser(ctx *fiber.Ctx) error {
	if os.Getenv("ADMIN_KEY") != ctx.Query("key") {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"VALIDATION_ERROR"},
			"message": "You are missing a API key.",
		})
	}

	status := util.BanUser(ctx.Query("email"), os.Getenv("ADMIN_KEY"))

	if status != 200 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"UNKNOWN_ERROR"},
			"message": "An unknown error occurred",
		})
	}

	col := structs.DB.Collection("users")

	var s structs.User

	if err := col.FindOneAndUpdate(context.TODO(), bson.M{"email": ctx.Query("email")}, bson.M{"$set": bson.M{
		"blacklisted": bson.M{
			"by":            ctx.Query("by"),
			"reason":        ctx.Query("reason"),
			"isblacklisted": true,
		},
	}}).Decode(&s); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"UNKNOWN_ERROR"},
			"message": "There was an error updating the blacklisted status for the user",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "OK",
		"errors":  structs.Errors{},
		"message": "Updated user successfully.",
	})
}
