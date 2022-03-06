package routes

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"shiba-backend/structs"
	"shiba-backend/util"
	"strings"
)

func complete(ctx *fiber.Ctx) error {
	Invite := ctx.Query("invite")
	Email := ctx.Query("email")
	Password := util.String(10)

	col := structs.DB.Collection("invites")

	var e bson.M

	if err := col.FindOne(context.TODO(), bson.M{
		"invite": Invite,
		"active": true,
	}).Decode(&e); err != nil {
		return ctx.Status(403).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  errors{"INVALID_INVITE"},
			"message": "The invite you provided was invalid. Please try again later.",
		})
	}

	if _, err := col.UpdateOne(context.TODO(), bson.M{"invite": Invite}, bson.M{"$set": bson.M{"active": false}}); err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  errors{"UNEXPECTED_ERROR"},
			"message": "An unexpected error has occurred",
			"error":   err.Error(),
		})
	}

	status := util.BasicSender(os.Getenv("API_URL")+"/admin/mail/users/add", Email, Password, os.Getenv("ADMIN_KEY"))

	if status != 200 {
		return ctx.Status(403).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  errors{"UNEXPECTED_ERROR"},
			"message": "An unexpected error has occurred",
		})
	}

	return ctx.JSON(fiber.Map{
		"status":  "OK",
		"errors":  errors{},
		"message": "Registered account successfully",
		"user": fiber.Map{
			"email":    Email,
			"invite":   Invite,
			"password": Password, // because this is auto generated serverside we have to provide this
		},
	})
}

func Register(ctx *fiber.Ctx) error {

	//Invite := ctx.Query("invite")
	Email := ctx.Query("email")
	//Password := util.String(10)

	if strings.Contains(Email, "@") == false {
		return ctx.Status(403).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  errors{"VALIDATION_ERROR"},
			"message": "You are missing an @ in your email.",
		})
	}

	for _, v := range util.GetDomains() {
		if strings.Contains(Email, v) == true {
			return complete(ctx)
		}

		return ctx.Status(403).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  errors{"VALIDATION_ERROR"},
			"message": "You are missing a valid domain in your email.",
		})
	}

	return nil
}
