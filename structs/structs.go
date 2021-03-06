package structs

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var DB *mongo.Database

type Errors []string

type StatsResponse struct {
	Domain string `json:"domain"`
	Users  []struct {
		Email      string   `json:"email"`
		Privileges []string `json:"privileges"`
		Status     string   `json:"status"`
		Mailbox    string   `json:"mailbox"`
	} `json:"users"`
}

type AliasesResponse []struct {
	Domain  string `json:"domain"`
	Aliases []struct {
		Address          string   `json:"address"`
		AddressDisplay   string   `json:"address_display"`
		ForwardsTo       []string `json:"forwards_to"`
		PermittedSenders []string `json:"permitted_senders"`
		Required         bool     `json:"required"`
	} `json:"aliases"`
}

type DomainsSlice []string

type Invite struct {
	Invite string `json:"invite"`
	MadeBy string `json:"madeby"`
	Active bool   `json:"active"`
	UsedBy struct {
		Email string `json:"email"`
		Date  int64  `json:"date"`
	} `json:"usedBy"`
}

// This is a constant that is used by main.go to show if the application is in testing process.
// If it is, all admin routes will not require a key
var IsTesting bool

type User struct {
	Email       string `json:"email"`
	Blacklisted struct {
		By          string `json:"by"`
		Reason      string `json:"reason"`
		Blacklisted bool   `json:"blacklisted"`
	} `json:"blacklisted"`
	Invite struct {
		MadeBy string `json:"madeby"`
		Date   int64  `json:"date"`
		Used   string `json:"used"`
	} `json:"invite"`
	Discord struct {
		ID int64 `json:"id"`
	} `json:"discord"`
}
