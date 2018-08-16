package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// collections names
const (
	AccountsCollection      = "accounts"
	TokenManagersCollection = "token_managers"
	RecruitsCollection      = "recruits"
	IndustriesCollection    = "industries"
)

// SetupEnv ...
func SetupEnv() {
	env := string(os.Getenv("ENV"))
	if env == "" {
		log.Println("Environment variable `ENV` not set, falling back to default config.")
		return
	}
	if err := godotenv.Load(".env." + env); err != nil {
		log.Println("Failed to load env file, falling back to default config.")
		os.Setenv("PORT", "9999")
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_NAME", "irecruit")
	}
}
