package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	// Load .env hanya jika di lingkungan lokal (development)
	if os.Getenv("GIN_MODE") != "release" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Tidak dapat memuat .env file, menggunakan environment variables yang ada.")
		}
	}
}
