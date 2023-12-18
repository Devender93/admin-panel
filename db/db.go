package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/kominkamen/rootds-admin/models"
)

var DB *pgxpool.Pool

func DBConfigFromEnv() *models.DBConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error In Connection")
	}

	return &models.DBConfig{
		DB_User:     os.Getenv("DB_USER"),
		DB_Pass:     os.Getenv("DB_PASS"),
		DB_Host:     os.Getenv("DB_HOST"),
		DB_Port:     os.Getenv("DB_PORT"),
		DB_Database: os.Getenv("DB_DATABASE"),
		DB_Sslmode:  os.Getenv("DB_SSLMODE"),
	}
}

func ConnectToDB() (*pgxpool.Pool, error) {

	dbConfig := DBConfigFromEnv()
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbConfig.DB_User,
		dbConfig.DB_Pass,
		dbConfig.DB_Host,
		dbConfig.DB_Port,
		dbConfig.DB_Database,
		dbConfig.DB_Sslmode,
	)

	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	return pool, nil
}
