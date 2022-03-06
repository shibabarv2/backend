package config

import (
	"github.com/joho/godotenv"
	"os"
)

func Get(key string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}

	str := os.Getenv("key")

	return str, nil

}
