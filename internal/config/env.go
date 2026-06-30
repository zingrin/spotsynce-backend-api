package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	ENV          string
	PORT         string
	DSN          string
	JWT_SECRET   string
	FRONTEND_URL string
}

func LoadEnv() *Env {
	err := godotenv.Load()

	if err != nil {
		// panic("Failed to load env file")
		fmt.Println("No .env file found, using system environment variables")
	}

	return &Env{
		ENV:          os.Getenv("ENV"),
		PORT:         os.Getenv("PORT"),
		DSN:          os.Getenv("DSN"),
		JWT_SECRET:   os.Getenv("JWT_SECRET"),
		FRONTEND_URL: os.Getenv("FRONTEND_URL"),
	}
}
