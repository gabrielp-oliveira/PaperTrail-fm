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

	createRootProjectTable := `
CREATE TABLE IF NOT EXISTS rootpappers (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	user_id TEXT,
	FOREIGN KEY(user_id) REFERENCES users(id)
);
`
	createPappersTable := `
CREATE TABLE IF NOT EXISTS pappers (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	path TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	root_papper_id TEXT,
	FOREIGN KEY(root_papper_id) REFERENCES rootpappers(id)
);
`

	CreateChapterTable := `
CREATE TABLE IF NOT EXISTS chapters (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	Papper_id TEXT,
	FOREIGN KEY(Papper_id) REFERENCES pappers(id)
);
`

	_, err := db.Exec(createRootProjectTable)
	if err != nil {
		log.Fatalf("Could not create Rootpappers table: %v", err)
	}

	_, err = db.Exec(createPappersTable)
	if err != nil {
		log.Fatalf("Could not create pappers table: %v", err)
	}
	_, err = db.Exec(CreateChapterTable)
	if err != nil {
		log.Fatalf("Could not create chapter table: %v", err)
	}

}
