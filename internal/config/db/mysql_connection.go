package dbconfig

import (
	config "github.com/ilhamfi27/ddd-golang-template/internal/config/app"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func (conn *DBConfig) NewMySQLConnection() (*gorm.DB, error) {
	env := config.GetEnvVars()
	dsn := env.DBUser + ":" + env.DBPass + "@tcp(" + env.DBHost + ":" + env.DBPort + ")/" + env.DBName + "?parseTime=true"
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
