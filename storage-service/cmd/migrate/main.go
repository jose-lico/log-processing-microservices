package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "postgres"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database: ", err)
	}

	fmt.Println("Successfully connected to the database!")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS process_log_entries (
			id SERIAL PRIMARY KEY,
			timestamp TIMESTAMP NOT NULL,
			level VARCHAR(10) NOT NULL CHECK (level IN ('INFO', 'WARN', 'ERROR', 'DEBUG')),
			message TEXT NOT NULL,
			user_id UUID,
			additional_data JSONB,
			processed BOOLEAN NOT NULL DEFAULT false
		)
	`)
	if err != nil {
		log.Fatal("Error creating table: ", err)
	}

	fmt.Println("Table 'process_log_entries' created successfully!")
}
