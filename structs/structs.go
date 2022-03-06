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
