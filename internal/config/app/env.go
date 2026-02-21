package appconfig

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVars struct {
	ServiceName    string
	DBDriver       string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPass         string
	DBName         string
	Port           string
	MigrationDir   string
	JWTSecret      string
	JaegerEndpoint string
	AppEnv         string
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

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
	}

	jaegerEndpoint := os.Getenv("JAEGER_ENDPOINT")
	if jaegerEndpoint == "" {
		jaegerEndpoint = "localhost:4318"
	}

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "ddd-golang-template"
	}

	return &EnvVars{
		ServiceName:    serviceName,
		DBDriver:       dBDriver,
		DBHost:         dBHost,
		DBPort:         dBPort,
		DBUser:         dBUser,
		DBPass:         dBPass,
		DBName:         dBName,
		Port:           ":" + port,
		MigrationDir:   migrationDir,
		JWTSecret:      jwtSecret,
		JaegerEndpoint: jaegerEndpoint,
		AppEnv:         appEnv,
	}
}
