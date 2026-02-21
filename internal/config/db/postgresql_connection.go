package dbconfig

import (
	config "github.com/ilhamfi27/ddd-golang-template/internal/config/app"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func (conn *DBConfig) NewPostgreSQLConnection() (*gorm.DB, error) {
	env := config.GetEnvVars()
	dsn := "host=" + env.DBHost + " user=" + env.DBUser + " password=" + env.DBPass + " dbname=" + env.DBName + " port=" + env.DBPort + " sslmode=disable"
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
