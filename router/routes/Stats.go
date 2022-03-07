package routes

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"os"
	"shiba-backend/structs"
	"shiba-backend/util"
)

func Stats(ctx *fiber.Ctx) error {
	key := os.Getenv("ADMIN_KEY")

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

		if util.Contains(domains, v.Domain) == false {
			domains = append(domains, v.Domain)
		}
	}

	aliasesT, err := http.NewRequest("GET", os.Getenv("API_URL")+"/admin/mail/aliases?format=json", nil)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"UNKNOWN_ERROR"},
			"message": "An unknown error has occurred while attempting to fetch stats. Please contact the developers of this application and wait until they fix it.",
			"error":   err.Error(),
		})
	}

	aliasesT.Header.Add("Authorization", "Basic "+key)

	reqz, err := client.Do(aliasesT)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "ERROR",
			"errors":  structs.Errors{"UNKNOWN_ERROR"},
			"message": "An unknown error has occurred while attempting to fetch stats. Please contact the developers of this application and wait until they fix it.",
			"error":   err.Error(),
		})
	}

	defer reqz.Body.Close()

	var x structs.AliasesResponse
	err = json.NewDecoder(reqz.Body).Decode(&x)

	var aliases int

	for _, _ = range x {
		aliases = aliases + 1
	}

	return ctx.JSON(fiber.Map{
		"status":  "OK",
		"errors":  structs.Errors{},
		"message": "Here is a list of statistics.",
		"stats": fiber.Map{
			"users":       users,
			"domains":     domains,
			"blacklisted": blacklisted,
			"aliases":     aliases,
		},
	})
}
