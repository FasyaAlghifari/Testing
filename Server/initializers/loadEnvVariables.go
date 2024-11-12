package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	// Jika Railway sudah mengatur environment, kita tidak perlu memuat .env
	if os.Getenv("RAILWAY_ENVIRONMENT") != "" {
		log.Println("Environment variables are already set by Railway.")
		return
	}

	// Jika tidak di Railway, muat dari file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, continuing with system env: ", err)
	}
}
