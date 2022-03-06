package routes

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"math/rand"
	"os"
	"shiba-backend/structs"
	"shiba-backend/util"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}

func Register(ctx *fiber.Ctx) error {

	Invite := ctx.Query("invite")
	Email := ctx.Query("email")
	Password := String(10)

	if strings.Contains(Email, "@") == false {
		return ctx.JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  errors{"VALIDATION_ERROR"},
			"message": "You are missing an @ in your email.",
		})
	}

	col := structs.DB.Collection("invites")

	var e bson.M

	if err := col.FindOne(context.TODO(), bson.M{
		"invite": Invite,
		"active": true,
	}).Decode(&e); err != nil {
		return ctx.JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  errors{"INVALID_INVITE"},
			"message": "The invite you provided was invalid. Please try again later.",
		})
	}

	if _, err := col.UpdateOne(context.TODO(), bson.M{"invite": Invite}, bson.M{"$set": bson.M{"active": false}}); err != nil {
		return ctx.JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  errors{"UNEXPECTED_ERROR"},
			"message": "An unexpected error has occurred",
			"error":   err.Error(),
		})
	}

	status := util.BasicSender("https://mail.shiba.bar/admin/mail/users/add", Email, Password, os.Getenv("ADMIN_KEY"))

	if status != 200 {
		return ctx.JSON(fiber.Map{
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
