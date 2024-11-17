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

		`CREATE TABLE IF NOT EXISTS Papers (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			path TEXT NOT NULL,
			"order" integer,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			world_id TEXT,
			FOREIGN KEY(world_id) REFERENCES worlds(id)
		);`,

		`CREATE TABLE IF NOT EXISTS timelines (
			id TEXT PRIMARY KEY,
			name TEXT,
			description TEXT,
			"order" integer,
			range integer,
			start_range integer,
			world_id TEXT NOT NULL,
			FOREIGN KEY(world_id) REFERENCES worlds(id)
		);`,

		`CREATE TABLE IF NOT EXISTS events (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			start_range integer NOT NULL,
			end_range integer NOT NULL,
			world_id TEXT,
			FOREIGN KEY(world_id) REFERENCES worlds(id)
		);`,

		`CREATE TABLE IF NOT EXISTS chapters (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			Paper_id TEXT,
			world_id TEXT,
			storyline_id TEXT,
			event_id TEXT,
			"order" integer,
			range integer DEFAULT 1,
			last_update DATE, 
			update TEXT,
			FOREIGN KEY (Paper_id) REFERENCES Papers(id),
			FOREIGN KEY (world_id) REFERENCES worlds(id),
			FOREIGN KEY (storyline_id) REFERENCES storyLines(id) ON DELETE SET NULL,
			FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE SET NULL
		);
	`,
		`CREATE TABLE IF NOT EXISTS storyLines (
			id TEXT PRIMARY KEY,
			name TEXT,
			description TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			world_id TEXT,
			"order" integer,
			FOREIGN KEY (world_id) REFERENCES worlds(id)
		);
	`,

		`CREATE TABLE IF NOT EXISTS connections (
			id TEXT PRIMARY KEY,
			source_chapter_id TEXT NOT NULL,
			target_chapter_id TEXT NOT NULL,
			world_id TEXT,
    		group_id TEXT DEFAULT '',
			FOREIGN KEY (source_chapter_id) REFERENCES chapters(id),
			FOREIGN KEY (target_chapter_id) REFERENCES chapters(id),
			FOREIGN KEY (world_id) REFERENCES worlds(id),
		);`,
		`CREATE TABLE IF NOT EXISTS group_connections (
			id TEXT PRIMARY KEY,
			world_id TEXT,
			color TEXT,
			name TEXT,
			description TEXT,
			FOREIGN KEY (world_id) REFERENCES worlds(id)
		);
`,
		`CREATE TABLE IF NOT EXISTS subway_settings (
			id TEXT PRIMARY KEY,
			chapter_names BOOLEAN,
			display_table_chapters BOOLEAN,
			timeline_update_chapter BOOLEAN,
			storyline_update_chapter BOOLEAN,
			zoom REAL,       
			x REAL,          
			y REAL,          
			world_id TEXT,
			FOREIGN KEY (world_id) REFERENCES worlds(id)
		);

		`,
	}

	for _, query := range createTables {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatal(" error -> ", query, err)
		}
	}

	fmt.Println("Tabelas criadas com sucesso!")
}
