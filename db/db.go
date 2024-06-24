package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "myuser"     // Substitua pelo seu usuário
	password = "mypassword" // Substitua pela sua senha
	dbname   = "mydatabase" // Substitua pelo seu banco de dados
)

var DB *sql.DB

func InitDB() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	createTables(DB)
}

func createTables(db *sql.DB) {
	createUsersTable := `
CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	email TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL
);
`

	_, err := db.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("Could not create users table: %v", err)
	}

	createEventsTable := `
CREATE TABLE IF NOT EXISTS events (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	location TEXT NOT NULL,
	dateTime TIMESTAMP NOT NULL,
	user_id INTEGER,
	FOREIGN KEY(user_id) REFERENCES users(id)
);
`

	_, err = db.Exec(createEventsTable)
	if err != nil {
		log.Fatalf("Could not create events table: %v", err)
	}

	createRegistrationsTable := `
CREATE TABLE IF NOT EXISTS registrations (
	id SERIAL PRIMARY KEY,
	event_id INTEGER,
	user_id INTEGER,
	FOREIGN KEY(event_id) REFERENCES events(id),
	FOREIGN KEY(user_id) REFERENCES users(id)
);
`

	_, err = db.Exec(createRegistrationsTable)
	if err != nil {
		log.Fatalf("Could not create registrations table: %v", err)
	}
}
