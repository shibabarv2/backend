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
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"VALIDATION_ERROR"},
			"message": "You are missing a API key.",
		})
	}

	if strings.HasPrefix(ctx.Query("invite"), "SHIB-") == false {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"VALIDATION_ERROR"},
			"message": "Valid invites start with SHIB. Please make your invite start with SHIB and then try again.",
		})
	}

	col := structs.DB.Collection("invites")

	if _, err := col.InsertOne(context.TODO(), bson.M{"invite": ctx.Query("invite"), "active": true, "usedBy": bson.M{"email": "", "date": "Never"}}); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"UNEXPECTED_ERROR"},
			"message": "An error occurred adding your invite.",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "OK",
		"errors":  structs.Errors{},
		"message": "Your invite has been successfully added to the database.",
		"invite": fiber.Map{
			"invite": ctx.Query("invite"),
			"active": true,
		},
	})
}
