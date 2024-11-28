package models

import (
	"database/sql"
	"fmt"

	"PaperTrail-fm.com/db"
)

type Event struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Range       int    `json:"range"`
	StartRange  int    `json:"startRange"`
	World_id    string `json:"world_id"`
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
			INSERT INTO events (id, name, description, range, start_range, world_id)
			VALUES ($1, $2, $3, $4, $5, $6)
		`
		_, err := db.DB.Exec(query, ev.Id, ev.Name, ev.Description, ev.Range, ev.StartRange, ev.World_id)
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

func (ev *Event) Update() error {

	// 	query := `
	// 	INSERT INTO events (id, name, description, range, start_range, world_id)
	// 	VALUES ($1, $2, $3, $4, $5, $6)
	// `

	query := `
	UPDATE events 
	SET name = $1, description = $2, range = $3, start_range = $4 WHERE id = $5

	`
	_, err := db.DB.Exec(query, ev.Name, ev.Description, ev.Range, ev.StartRange, ev.Id)
	if err != nil {
		return fmt.Errorf("error updating event: %v", err)
	}

	fmt.Println("Updated event in database")

	return nil
}

func getEventsByWorldId(worldId string) ([]Event, error) {
	events := []Event{}

	eventsQuery := `
	SELECT id, name, range, start_range, world_id
	FROM events
	WHERE world_id = $1
`
	rows, err := db.DB.Query(eventsQuery, worldId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event Event
		if err := rows.Scan(&event.Id, &event.Name, &event.Range, &event.StartRange, &event.World_id); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, err
}
