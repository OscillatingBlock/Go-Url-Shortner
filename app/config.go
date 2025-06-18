package app

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv" // Import godotenv
)

type Configuration struct {
	BaseUrl string
	Port    string
}

func (c Configuration) String() string {
	return fmt.Sprintf("BASE_URL=%s", c.BaseUrl)
}

func InitConfig() Configuration {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Warning: .env file not found or could not be loaded. Relying on system environment variables: ", "error", err)
		// Don't fatal here, allow fallback to system env vars or default config
	} else {
		slog.Error(".env file loaded successfully.")
	}
	base_url := os.Getenv("BASE_URL")
	if base_url == "" {
		slog.Warn("BASE_URL env not found, defaulting to http://localhost:8080/")
		base_url = "http://localhost:8080/"
	}

	port := os.Getenv("PORT")
	if port == "" {
		slog.Warn("PORT env not found, defaulting to 8080")
		port = "8080"
	}

	c := Configuration{BaseUrl: base_url, Port: port}
	slog.Info("config:")
	slog.Info(c.String())
	return c
}
