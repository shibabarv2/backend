package routes

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"shiba-backend/structs"
	"strings"
	"time"
)

func GetInvite(ctx *fiber.Ctx) error {
	if os.Getenv("ADMIN_KEY") != ctx.Query("key") {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"VALIDATION_ERROR"},
			"message": "You are missing an API key.",
		})
	}

	if strings.HasPrefix(ctx.Query("invite"), "SHIB-") == false {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"VALIDATION_ERROR"},
			"message": "Please ensure the invite starts with SHIB and is valid before trying again",
		})
	}

	col := structs.DB.Collection("invites")

	var a structs.Invite

	if err := col.FindOne(context.TODO(), bson.M{"invite": ctx.Query("invite")}).Decode(&a); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"UNEXPECTED_ERROR"},
			"message": "An error occurred looking up the invite.",
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "OK",
		"errors":  structs.Errors{},
		"message": "Here is the invite information.",
		"invite": fiber.Map{
			"invite": ctx.Query("invite"),
			"active": a.Active,
			"madeBy": a.MadeBy,
			"usedBy": fiber.Map{
				"email": a.UsedBy.Email,
				"date":  time.Unix(a.UsedBy.Date, 0),
				"unix":  a.UsedBy.Date,
			},
		},
	})
}
