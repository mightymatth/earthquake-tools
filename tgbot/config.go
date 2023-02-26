package main

import (
	"github.com/joho/godotenv"
	"os"
)

// TODO: use viper (with cobra) instead of this
func getParam(param, envKey, defaultVal string) string {
	if param == "" {
		return getEnv(envKey, defaultVal)
	}

	return param
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func loadDotEnv() error {
	return godotenv.Load()
}
