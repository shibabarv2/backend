package auth

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"shiba-backend/structs"
	"shiba-backend/util"
	"strings"
	"time"
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
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"INVALID_INVITE"},
			"message": "The invite you provided was invalid. Please try again later.",
		})
	}

	if _, err := col.UpdateOne(context.TODO(), bson.M{"invite": Invite}, bson.M{"$set": bson.M{"active": false, "usedBy": bson.M{"email": Email, "date": time.Now().String()}}}); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"UNEXPECTED_ERROR"},
			"message": "An unexpected error has occurred",
			"error":   err.Error(),
		})
	}

	status := util.BasicSender(os.Getenv("API_URL")+"/admin/mail/users/add", Email, Password, os.Getenv("ADMIN_KEY"))

	if status != 200 {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"UNEXPECTED_ERROR"},
			"message": "An unexpected error has occurred",
		})
	}

	return ctx.JSON(fiber.Map{
		"status":  "OK",
		"errors":  structs.Errors{},
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
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"VALIDATION_ERROR"},
			"message": "You are missing an @ in your email.",
		})
	}

	for _, v := range util.GetDomains() {
		if strings.Contains(Email, v) == true {
			return complete(ctx)
		}
	}

	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"status":  "ERROR",
		"errors":  structs.Errors{"VALIDATION_ERROR"},
		"message": "You are missing a valid domain in your email.",
	})
}
