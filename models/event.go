package models

import (
	"database/sql"
	"fmt"

	"PaperTrail-fm.com/db"
)

type Event struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	Start_date     string `json:"start_date"`
	End_date       string `json:"end_date"`
	Root_papper_id string `json:"root_papper_id"`
}

func (ev *Event) Save() error {
	var id string

	// Verifica se o evento já existe no banco de dados
	query := `SELECT id FROM events WHERE id = $1`
	err := db.DB.QueryRow(query, ev.Id).Scan(&id)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking event existence: %v", err)
	}

	if err == sql.ErrNoRows {
		// Inserir novo evento na tabela events
		query := `
			INSERT INTO events (id, name, start_date, end_date, root_papper_id)
			VALUES ($1, $2, $3, $4, $5)
		`
		_, err := db.DB.Exec(query, ev.Id, ev.Name, ev.Start_date, ev.End_date, ev.Root_papper_id)
		if err != nil {
			return fmt.Errorf("error inserting event: %v", err)
		}

		fmt.Println("Inserted new event into database")
	} else {
		fmt.Println("Event already exists in database")
	}

	return nil
}

func (ev *Event) Delete() error {
	query := `SELECT id FROM events WHERE id = $1`
	err := db.DB.QueryRow(query, ev.Id).Scan(&ev.Id)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking event existence: %v", err)
	}
	query = `DELETE FROM events WHERE id = $1;`
	_, err = db.DB.Exec(query, ev.Id)
	if err != nil {
		return fmt.Errorf("error removing event: %v", err)
	}

	fmt.Println("Removed event from database")

	return nil
}
