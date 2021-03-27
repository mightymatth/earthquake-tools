package main

import (
	"github.com/joho/godotenv"
	"os"
)

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func loadDotEnv() error {
	return godotenv.Load()
}
