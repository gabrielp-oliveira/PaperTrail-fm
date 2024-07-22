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

	createTables := []string{

		`CREATE TABLE IF NOT EXISTS rootpappers (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			user_id TEXT,
			FOREIGN KEY(user_id) REFERENCES users(id)
		);`,

		`CREATE TABLE IF NOT EXISTS pappers (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			path TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			root_papper_id TEXT,
			FOREIGN KEY(root_papper_id) REFERENCES rootpappers(id)
		);`,

		`CREATE TABLE IF NOT EXISTS timelines (
			id TEXT PRIMARY KEY,
			root_papper_id TEXT NOT NULL,
			date DATE NOT NULL,
			FOREIGN KEY(root_papper_id) REFERENCES rootpappers(id)
		);`,

		`CREATE TABLE IF NOT EXISTS events (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			start_date DATE NOT NULL,
			end_date DATE NOT NULL,
			root_papper_id TEXT,
			FOREIGN KEY(root_papper_id) REFERENCES rootpappers(id)
		);`,

		`CREATE TABLE IF NOT EXISTS chapters (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			papper_id TEXT,
			root_papper_id TEXT,
			time_line_id TEXT,
			event_id TEXT,
			FOREIGN KEY(papper_id) REFERENCES pappers(id),
			FOREIGN KEY(root_papper_id) REFERENCES rootpappers(id),
			FOREIGN KEY(time_line_id) REFERENCES timelines(id),
			FOREIGN KEY(event_id) REFERENCES events(id)
		);`,

		`CREATE TABLE IF NOT EXISTS connections (
			id TEXT PRIMARY KEY,
			source_chapter_id TEXT NOT NULL,
			target_chapter_id TEXT NOT NULL,
			FOREIGN KEY(source_chapter_id) REFERENCES chapters(id),
			FOREIGN KEY(target_chapter_id) REFERENCES chapters(id)
		);`,
	}

	for _, query := range createTables {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Tabelas criadas com sucesso!")
}
