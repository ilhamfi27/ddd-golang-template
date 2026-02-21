package appconfig

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVars struct {
	DBDriver     string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPass       string
	DBName       string
	Port         string
	MigrationDir string
}

func LoadEnv() error {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "production" {
		return nil
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return nil
}

func GetEnvVars() *EnvVars {
	dBDriver := os.Getenv("DATABASE_DRIVER")
	dBHost := os.Getenv("DATABASE_HOST")
	dBPort := os.Getenv("DATABASE_PORT")
	dBUser := os.Getenv("DATABASE_USER")
	dBPass := os.Getenv("DATABASE_PASS")
	dBName := os.Getenv("DATABASE_NAME")
	migrationDir := os.Getenv("MIGRATION_DIR")
	if migrationDir == "" {
		migrationDir = "./internal/config/db/migrations"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "1321"
	}
	return &EnvVars{
		DBDriver:     dBDriver,
		DBHost:       dBHost,
		DBPort:       dBPort,
		DBUser:       dBUser,
		DBPass:       dBPass,
		DBName:       dBName,
		Port:         ":" + port,
		MigrationDir: migrationDir,
	}
}
