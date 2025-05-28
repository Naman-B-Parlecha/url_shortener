package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	godotenv.Load()
	val := os.Getenv("DSN")
	connectionStr := val

	conn, err := sql.Open("postgres", connectionStr)
	if err != nil {
		log.Println("Failed to connect to the database:", err)
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		log.Println("Failed to Ping the database:", err)
		return nil, err
	}
	log.Println("Connected to the database successfully")
	return conn, nil
}
