// Package config
package config

import (
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SecretKey           string
	DBConnectionString  string
	LogFormat           string
	LogLevel            string
	RedisClientHost     string
	RedisClientPassword string
}

func (c *Config) Load() {
	if err := godotenv.Load(); err != nil {
		slog.Info("godot env load failed", "error", err)
	}

	c.SecretKey = os.Getenv("SECRET_KEY")
	if c.SecretKey == "" {
		log.Fatal("secret key is empty!")
	}
	c.DBConnectionString = os.Getenv("DB_CONNECTION_STRING")
	if c.DBConnectionString == "" {
		log.Fatal("db connection string is empty!")
	}
	c.LogFormat = os.Getenv("LOG_FORMAT")
	c.LogLevel = os.Getenv("LOG_LEVEL")
	c.RedisClientHost = os.Getenv("REDIS_CLIENT_HOST")
	c.RedisClientPassword = os.Getenv("REDIS_CLIENT_PASSWORD")
	if c.RedisClientPassword == "" {
		log.Fatal("redis password is empty!")
	}
}
