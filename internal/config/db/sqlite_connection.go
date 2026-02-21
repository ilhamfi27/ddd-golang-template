package dbconfig

import (
	"log"
	"os"
	"path/filepath"

	config "github.com/ilhamfi27/ddd-golang-template/internal/config/app"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func (conn *DBConfig) NewSqliteConnection() (*gorm.DB, error) {
	newPath := filepath.Join("dbs")
	err := os.MkdirAll(newPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	return gorm.Open(sqlite.Open(config.SQLITE_LOCATION), &gorm.Config{})
}
