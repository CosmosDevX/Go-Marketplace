package helpers

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

type TestsConfig struct {
	SecretKey               string
	TestsDBConnectionString string
}

func LoadTestConfig() TestsConfig {
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(filename), "../..")
	envPath := filepath.Join(root, ".env")

	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("failed to load .env for tests: %v", err)
	}

	return TestsConfig{
		SecretKey:               os.Getenv("SECRET_KEY"),
		TestsDBConnectionString: os.Getenv("TESTS_DB_CONNECTION_STRING"),
	}
}
