package auth

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"shiba-backend/structs"
	"shiba-backend/util"
	"strconv"
	"strings"
	"time"
)

func complete(ctx *fiber.Ctx) error {
	Invite := ctx.Query("invite")
	Email := ctx.Query("email")
	DiscordID := ctx.Query("id")
	Password := util.String(10)

	col := structs.DB.Collection("invites")
	users := structs.DB.Collection("users")

	var e bson.M

	if err := col.FindOne(context.TODO(), bson.M{
		"invite": Invite,
		"active": true,
	}).Decode(&e); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"INVALID_INVITE"},
			"message": "The invite provided was invalid. Please try again later.",
		})
	}

	if _, err := col.UpdateOne(context.TODO(), bson.M{"invite": Invite}, bson.M{"$set": bson.M{"active": false, "usedBy": bson.M{"email": Email, "date": time.Now().Unix()}}}); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"UNEXPECTED_ERROR"},
			"message": "An unexpected error has occurred.",
			"error":   err.Error(),
		})
	}

	status, _ := util.BasicSender(os.Getenv("API_URL")+"/admin/mail/users/add", Email, Password, os.Getenv("ADMIN_KEY"))

	if status != 200 {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"UNEXPECTED_ERROR"},
			"message": "An unexpected error has occurred.",
		})
	}

	// Attempt to make the DiscordID a int64
	i, err := strconv.ParseInt(DiscordID, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"UNEXPECTED_ERROR"},
			"message": "An unexpected error has occurred.",
		})
	}

	if _, err := users.InsertOne(context.TODO(), bson.M{
		"email": Email,
		"blacklisted": bson.M{
			"reason":      nil,
			"by":          nil,
			"blacklisted": false,
		}, "invite": bson.M{
			"madeby": e["madeby"],
			"date":   time.Now().Unix(),
			"used":   Invite,
		},
		"discord": bson.M{
			"id": i,
		},
	}); err != nil {
		return ctx.JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"UNKNOWN_ERROR"},
			"message": "An error has occurred and your account has not been registered.",
			"error":   err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"status":  "OK",
		"errors":  structs.Errors{},
		"message": "Registered account successfully.",
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
			"message": "The email provided is missing an \"@\".",
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
		"message": "The email provided is missing a valid domain.",
	})
}
