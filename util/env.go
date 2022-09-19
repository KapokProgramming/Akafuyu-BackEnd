package util

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func MustGet(key string) string {
	val := os.Getenv(key)

	if val == "" {
		log.Println("Env key is missing " + key)
	}

	return val
}

func CheckENV() {
	err := godotenv.Load()

	if err != nil {
		log.Panicln("Error loading .env file")
	}
}
