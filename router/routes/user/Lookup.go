package user

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"os"
	"shiba-backend/structs"
)

func Lookup(ctx *fiber.Ctx) error {

	key := os.Getenv("ADMIN_KEY")

	client := &http.Client{}

	req, err := http.NewRequest("GET", os.Getenv("API_URL")+"/admin/mail/users?format=json", nil)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"UNKNOWN_ERROR"},
			"message": "An unknown error has occurred while attempting to fetch this user. Please contact the developers of this application and wait until they fix it.",
			"error":   err.Error(),
		})
	}

	req.Header.Add("Authorization", "Basic "+key)

	resp, err := client.Do(req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"UNKNOWN_ERROR"},
			"message": "An unknown error has occurred while attempting to fetch this user. Please contact the developers of this application and wait until they fix it.",
			"error":   err.Error(),
		})
	}

	defer resp.Body.Close()

	var r []structs.StatsResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"UNKNOWN_ERROR"},
			"message": "An unknown error has occurred while attempting to fetch this user. Please contact the developers of this application and wait until they fix it.",
			"error":   err.Error(),
		})
	}

	users := structs.DB.Collection("users")

	var User structs.User

	if err := users.FindOne(context.TODO(), bson.M{
		"email": ctx.Params("username"),
	}).Decode(&User); err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"USER_NOT_FOUND"},
			"message": "The user you queried could not be found.",
			"error":   err.Error(),
		})
	}

	if User.Email == ctx.Params("username") {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "OK",
			"errors": structs.Errors{},
			"user": fiber.Map{
				"email":       User.Email,
				"invite":      User.Invite,
				"blacklisted": User.Blacklisted,
				"discord":     User.Discord,
			},
		})
	}

	return nil

}
