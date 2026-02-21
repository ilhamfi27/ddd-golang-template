package handlers

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"

	config "github.com/ilhamfi27/ddd-golang-template/internal/config/app"
	dbconfig "github.com/ilhamfi27/ddd-golang-template/internal/config/db"
	"gorm.io/gorm"
)

type DBHandler struct {
}

func NewDBHandler() *DBHandler {
	return &DBHandler{}
}

func migrationDir() string {
	migrationDir := config.GetEnvVars().MigrationDir
	migrationDir = fmt.Sprintf("file://%s", migrationDir)
	return migrationDir
}

func handleMySQLConnection() *gorm.DB {
	db := dbconfig.NewDB()
	conn, err := db.NewMySQLConnection()
	if err != nil {
		panic(err)
	}
	return conn
}

func handlePostgresConnection() *gorm.DB {
	db := dbconfig.NewDB()
	conn, err := db.NewPostgreSQLConnection()
	if err != nil {
		panic(err)
	}
	return conn
}

func handleSqliteConnection() *gorm.DB {
	db := dbconfig.NewDB()
	conn, err := db.NewSqliteConnection()
	if err != nil {
		panic(err)
	}
	return conn
}

func handleMySQLMigration() error {
	env := config.GetEnvVars()
	dsn := env.DBUser + ":" + env.DBPass + "@tcp(" + env.DBHost + ":" + env.DBPort + ")/" + env.DBName + "?parseTime=true"
	conn, _ := sql.Open(dbconfig.MYSQL, dsn)
	driver, _ := mysql.WithInstance(conn, &mysql.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		migrationDir(),
		dbconfig.MYSQL,
		driver,
	)

	if err != nil {
		panic(err)
	}
	fmt.Println(m.Up())
	return nil
}

func handlePostgresMigration() error {
	env := config.GetEnvVars()
	dsn := "host=" + env.DBHost + " user=" + env.DBUser + " password=" + env.DBPass + " dbname=" + env.DBName + " port=" + env.DBPort + " sslmode=disable"
	conn, _ := sql.Open(dbconfig.POSTGRES, dsn)
	driver, _ := postgres.WithInstance(conn, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		migrationDir(),
		dbconfig.POSTGRES,
		driver,
	)

	if err != nil {
		panic(err)
	}
	fmt.Println(m.Up())
	return nil
}

func handleSqliteMigration() error {
	conn, _ := sql.Open(dbconfig.SQLITE, config.SQLITE_LOCATION)
	driver, _ := sqlite.WithInstance(conn, &sqlite.Config{})
	m, _ := migrate.NewWithDatabaseInstance(
		migrationDir(),
		dbconfig.SQLITE,
		driver,
	)

	return m.Up()
}

func (dbHandler *DBHandler) DBInit() *gorm.DB {
	dbDriver := config.GetEnvVars().DBDriver
	switch dbDriver {
	case dbconfig.MYSQL:
		return handleMySQLConnection()
	case dbconfig.POSTGRES:
		return handlePostgresConnection()
	case dbconfig.SQLITE:
		return handleSqliteConnection()
	default:
		return handleSqliteConnection()
	}
}

func (dbHandler *DBHandler) Migrate() error {
	dbDriver := config.GetEnvVars().DBDriver
	switch dbDriver {
	case dbconfig.MYSQL:
		return handleMySQLMigration()
	case dbconfig.POSTGRES:
		return handlePostgresMigration()
	case dbconfig.SQLITE:
		return handleSqliteMigration()
	default:
		return handleSqliteMigration()
	}
}
