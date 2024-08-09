package models

import (
	"database/sql"
	"fmt"
	"time"

	"PaperTrail-fm.com/db"
)

type Worlds struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UserID    string    `json:"user_id"`
}

type Env struct {
	Worlds
	Chapters    []Chapter    `json:"chapters"`
	Events      []Event      `json:"events"`
	Connections []Connection `json:"connections"`
	Timelines   []Timeline   `json:"timelines"`
}

// Salva o Worlds no banco de dados
func (rp *Worlds) Save() error {
	var rpId string

	query := `SELECT id FROM worlds WHERE id = $1`
	err := db.DB.QueryRow(query, rp.Id).Scan(&rpId)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking worlds existence: %v", err)
	}

	if err == sql.ErrNoRows {
		insertQuery := `
		INSERT INTO worlds(id, name, created_at, user_id) 
		VALUES ($1, $2, $3, $4)`
		_, err := db.DB.Exec(insertQuery, rp.Id, rp.Name, rp.CreatedAt, rp.UserID)
		if err != nil {
			return fmt.Errorf("error inserting world: %v", err)
		}

		fmt.Println("Inserted new world into database")
	} else {
		fmt.Println("World already exists in database")
	}

	return nil
}

// Obtém a lista de pappers associados ao Worlds
func (rp *Worlds) GetPapperList() ([]Papper, error) {
	query := "SELECT id, name, description, created_at FROM pappers WHERE world_id = $1"
	rows, err := db.DB.Query(query, rp.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Papper
	for rows.Next() {
		var papper Papper
		if err := rows.Scan(&papper.ID, &papper.Name, &papper.Description, &papper.Created_at); err != nil {
			return nil, err
		}
		list = append(list, papper)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

// Obtém capítulos, conexões, eventos e timelines associados ao Worlds
func (rp *Worlds) GetRootData() (Env, error) {
	chapters := []Chapter{}
	events := []Event{}
	connections := []Connection{}
	timelines := []Timeline{}
	env := Env{}

	// Consultar eventos
	eventsQuery := `
		SELECT id, name, start_date, end_date, world_id
		FROM events
		WHERE world_id = $1
	`
	rows, err := db.DB.Query(eventsQuery, rp.Id)
	if err != nil {
		return env, err
	}
	defer rows.Close()

	for rows.Next() {
		var event Event
		if err := rows.Scan(&event.Id, &event.Name, &event.Start_date, &event.End_date, &event.World_id); err != nil {
			return env, err
		}
		events = append(events, event)
	}

	// Consultar capítulos
	chaptersQuery := `
		SELECT id, name, description, created_at, papper_id, world_id, event_id
		FROM chapters
		WHERE world_id = $1
	`
	rows, err = db.DB.Query(chaptersQuery, rp.Id)
	if err != nil {
		return env, err
	}
	defer rows.Close()

	for rows.Next() {
		var chapter Chapter
		if err := rows.Scan(&chapter.Id, &chapter.Name, &chapter.Description, &chapter.CreatedAt,
			&chapter.PapperID, &chapter.WorldsID, &chapter.EventID); err != nil {
			return env, err
		}
		chapters = append(chapters, chapter)
	}

	// Consultar conexões
	connectionsQuery := `
		SELECT id, source_chapter_id, target_chapter_id
		FROM connections
		WHERE source_chapter_id IN (SELECT id FROM chapters WHERE world_id = $1)
	`
	rows, err = db.DB.Query(connectionsQuery, rp.Id)
	if err != nil {
		return env, err
	}
	defer rows.Close()

	for rows.Next() {
		var connection Connection
		if err := rows.Scan(&connection.Id, &connection.SourceChapterID, &connection.TargetChapterID); err != nil {
			return env, err
		}
		connections = append(connections, connection)
	}

	// Consultar timelines
	timelinesQuery := `
		SELECT id, name, date, world_id
		FROM timelines
		WHERE world_id = $1
	`
	rows, err = db.DB.Query(timelinesQuery, rp.Id)
	if err != nil {
		return env, err
	}
	defer rows.Close()

	for rows.Next() {
		var timeline Timeline
		if err := rows.Scan(&timeline.Id, &timeline.Date, &timeline.WorldsID); err != nil {
			return env, err
		}
		timelines = append(timelines, timeline)
	}

	// Preencher o objeto Env
	env.Worlds = *rp
	env.Chapters = chapters
	env.Connections = connections
	env.Events = events
	env.Timelines = timelines

	return env, nil
}
