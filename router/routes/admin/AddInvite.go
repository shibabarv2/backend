package routes

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"shiba-backend/structs"
	"strings"
)

func AddInvite(ctx *fiber.Ctx) error {
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

	if _, err := col.InsertOne(context.TODO(), bson.M{"invite": ctx.Query("invite"), "active": true, "madeby": "Admin (This user was invited by an administrator)", "usedBy": bson.M{"email": "", "date": "Never"}}); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"UNEXPECTED_ERROR"},
			"message": "An error occurred while adding the invite.",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "OK",
		"errors":  structs.Errors{},
		"message": "The invite has been successfully added to the database.",
		"invite": fiber.Map{
			"invite": ctx.Query("invite"),
			"active": true,
		},
	})
}
