package routes

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"shiba-backend/structs"
	"strings"
)

func GetInvite(ctx *fiber.Ctx) error {
	if os.Getenv("ADMIN_KEY") != ctx.Query("key") {
		return ctx.Status(403).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"VALIDATION_ERROR"},
			"message": "You are missing a API key.",
		})
	}

	if strings.HasPrefix(ctx.Query("invite"), "SHIB-") == false {
		return ctx.Status(403).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"VALIDATION_ERROR"},
			"message": "A valid invite starts with SHIB. Please lookup a valid invite and then try again",
		})
	}

	col := structs.DB.Collection("invites")

	var a structs.Invite

	if err := col.FindOne(context.TODO(), bson.M{"invite": ctx.Query("invite")}).Decode(&a); err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"UNEXPECTED_ERROR"},
			"message": "An error occurred looking up your invite.",
			"error":   err.Error(),
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "OK",
		"errors":  structs.Errors{},
		"message": "Here is the invite information.",
		"invite": fiber.Map{
			"invite": ctx.Query("invite"),
			"active": true,
			"usedBy": a.UsedBy,
		},
	})
}