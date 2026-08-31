// Package config
package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SecretKey               string
	DBConnectionString      string
	TestsDBConnectionString string
	LogFormat               string
	LogLevel                string
	RedisClientHost         string
	RedisClientPassword     string
	CORSAllowedHost         string
	UploadsDir              string
}

func (c *Config) Load() {
	if err := godotenv.Load(); err != nil {
		slog.Info("godot env load failed", "error", err)
	}

	c.SecretKey = os.Getenv("SECRET_KEY")
	if c.SecretKey == "" {
		slog.Error("secret key is empty")
		os.Exit(1)
	}
	c.DBConnectionString = os.Getenv("DB_CONNECTION_STRING")
	if c.DBConnectionString == "" {
		slog.Error("db connection string is empty")
		os.Exit(1)
	}
	c.LogFormat = os.Getenv("LOG_FORMAT")
	c.LogLevel = os.Getenv("LOG_LEVEL")
	c.RedisClientHost = os.Getenv("REDIS_CLIENT_HOST")
	c.RedisClientPassword = os.Getenv("REDIS_CLIENT_PASSWORD")
	c.CORSAllowedHost = os.Getenv("CORS_ALLOWED_HOST")
	c.UploadsDir = os.Getenv("UPLOADS_DIR")
	if c.UploadsDir == "" {
		slog.Error("uploads dir is empty")
		os.Exit(1)
	}
}
