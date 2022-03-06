package routes

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"os"
	"shiba-backend/structs"
)

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func Stats(ctx *fiber.Ctx) error {
	key := os.Getenv("ADMIN_KEY")

	client := &http.Client{}

	req, err := http.NewRequest("GET", os.Getenv("API_URL")+"/admin/mail/users?format=json", nil)

	if err != nil {
		return ctx.JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  errors{"UNKNOWN_ERROR"},
			"message": "An unknown error has occurred while attempting to fetch stats. Please contact the developers of this application and wait until they fix it.",
			"error":   err.Error(),
		})
	}

	req.Header.Add("Authorization", "Basic "+key)

	resp, err := client.Do(req)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  errors{"UNKNOWN_ERROR"},
			"message": "An unknown error has occurred while attempting to fetch stats. Please contact the developers of this application and wait until they fix it.",
			"error":   err.Error(),
		})
	}

	defer resp.Body.Close()

	var r structs.StatsResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  errors{"UNKNOWN_ERROR"},
			"message": "An unknown error has occurred while attempting to fetch stats. Please contact the developers of this application and wait until they fix it.",
			"error":   err.Error(),
		})
	}

	var users int
	var blacklisted int
	var domains []string

	users = 0

	for _, v := range r {
		for _, n := range v.Users {
			if n.Status == "inactive" {
				blacklisted = blacklisted + 1
			} else {
				users = users + 1
			}
		}

		if Contains(domains, v.Domain) == false {
			domains = append(domains, v.Domain)
		}
	}

	return ctx.JSON(fiber.Map{
		"status":      "OK",
		"users":       users,
		"blacklisted": blacklisted,
		"domains":     domains,
	})
}
