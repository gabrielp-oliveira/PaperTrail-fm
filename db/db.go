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

		`CREATE TABLE IF NOT EXISTS worlds (
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
			"order" integer,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			world_id TEXT,
			FOREIGN KEY(world_id) REFERENCES worlds(id)
		);`,

		`CREATE TABLE IF NOT EXISTS timelines (
			id TEXT PRIMARY KEY,
			world_id TEXT NOT NULL,
			name TEXT,
			date DATE NOT NULL,
			FOREIGN KEY(world_id) REFERENCES worlds(id)
		);`,

		`CREATE TABLE IF NOT EXISTS events (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			start_date DATE NOT NULL,
			end_date DATE NOT NULL,
			world_id TEXT,
			FOREIGN KEY(world_id) REFERENCES worlds(id)
		);`,

		`CREATE TABLE IF NOT EXISTS chapters (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			papper_id TEXT,
			world_id TEXT,
			timeline_id TEXT,
			event_id TEXT,
			"order" integer,
			last_update DATE, 
			update TEXT,
			FOREIGN KEY (papper_id) REFERENCES pappers(id),
			FOREIGN KEY (world_id) REFERENCES worlds(id),
			FOREIGN KEY (timeline_id) REFERENCES timelines(id) ON DELETE SET NULL,
			FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE SET NULL
		);
	`,

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
