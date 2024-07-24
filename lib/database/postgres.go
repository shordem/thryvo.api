package database

import (
	"fmt"
	"time"

	"github.com/shordem/api.thryvo/lib/constants"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DatabaseFacade *gorm.DB

type PostgresClientInterface interface {
	Connection() *gorm.DB
}

type pgClient struct {
	database *gorm.DB
}

func NewPostgresClient(env constants.Env) PostgresClientInterface {
	dsn := "host=" + env.DB_HOST + " user=" + env.DB_USER + " password=" + env.DB_PASSWORD + " dbname=" + env.DB_NAME + " port=" + env.DB_PORT

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

	// defer sqlDb.Close()

	DatabaseFacade = database

	return &pgClient{
		database: database,
	}
}

func (conn pgClient) Connection() *gorm.DB {
	return conn.database
}
