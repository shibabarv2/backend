package structs

import "go.mongodb.org/mongo-driver/mongo"

var DB *mongo.Database

type Errors []string

type StatsResponse []struct {
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
	Active bool   `json:"active"`
}
