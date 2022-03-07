package user

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"os"
	"shiba-backend/structs"
)

func Lookup(ctx *fiber.Ctx) error {

	key := os.Getenv("ADMIN_KEY")

	if key != ctx.Query("key") {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"VALIDATION_ERROR"},
			"message": "There was an error validating your key.",
		})
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", os.Getenv("API_URL")+"/admin/mail/users?format=json", nil)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"UNKNOWN_ERROR"},
			"message": "An unknown error has occurred while attempting to fetch stats. Please contact the developers of this application and wait until they fix it.",
			"error":   err.Error(),
		})
	}

	req.Header.Add("Authorization", "Basic "+key)

	resp, err := client.Do(req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"UNKNOWN_ERROR"},
			"message": "An unknown error has occurred while attempting to fetch stats. Please contact the developers of this application and wait until they fix it.",
			"error":   err.Error(),
		})
	}

	defer resp.Body.Close()

	var r structs.StatsResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"UNKNOWN_ERROR"},
			"message": "An unknown error has occurred while attempting to fetch stats. Please contact the developers of this application and wait until they fix it.",
			"error":   err.Error(),
		})
	}

	for _, v := range r {
		for _, n := range v.Users {
			if n.Email == ctx.Query("username") {
				return ctx.JSON(fiber.Map{
					"status":  "OK",
					"errors":  structs.Errors{},
					"message": "Here is the lookup info for the user",
					"user": fiber.Map{
						"email":      n.Email,
						"status":     n.Status,
						"privileges": n.Privileges,
						"mailbox":    n.Mailbox,
					},
				})
			}
		}

		return ctx.JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"USER_NOT_FOUND"},
			"message": "The user you queried could not be found.",
		})

	}

	return nil

}
