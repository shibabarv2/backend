package routes

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"shiba-backend/structs"
	"strings"
)

func RemoveInvite(ctx *fiber.Ctx) error {
	if os.Getenv("ADMIN_KEY") != ctx.Query("key") {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"VALIDATION_ERROR"},
			"message": "You are missing an API key.",
		})
	}

	if !strings.HasPrefix(ctx.Query("invite"), "SHIB-") {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"VALIDATION_ERROR"},
			"message": "Please ensure the invite starts with SHIB before trying again.",
		})
	}

	col := structs.DB.Collection("invites")

	var n structs.Invite

	if err := col.FindOne(context.TODO(), bson.M{"invite": ctx.Query("invite"), "active": true}).Decode(&n); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"INVALID_INVITE"},
			"message": "The invite provided was invalid.",
		})
	}

	if _, err := col.UpdateOne(context.TODO(), bson.M{"invite": ctx.Query("invite"), "active": true}, bson.M{"$set": bson.M{"active": false}}); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"UNEXPECTED_ERROR"},
			"message": "An error occurred removing the invite.",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "OK",
		"errors":  structs.Errors{},
		"message": "The invite has been successfully removed.",
	})
}
