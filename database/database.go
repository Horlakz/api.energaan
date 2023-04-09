package database

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DatabaseFacade *gorm.DB

type DatabaseInterface interface {
	Connection() *gorm.DB
}

type connection struct {
	database *gorm.DB
}

func StartDatabaseClient() DatabaseInterface {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true})

	if err != nil {
		fmt.Print(err)
	}

	sqlDb, err := database.DB()

	sqlDb.SetMaxIdleConns(5)
	sqlDb.SetMaxOpenConns(10)
	sqlDb.SetConnMaxLifetime(time.Hour)

	if err != nil {
		fmt.Print(err)
	}

	fmt.Println("Database connection is successful")

	database.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	database.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto")

	//defer db.Close()

	DatabaseFacade = database

	return &connection{
		database: database,
	}
}

func (conn connection) Connection() *gorm.DB {
	return conn.database
}
