package models

import (
	"database/sql"
	"fmt"
	"time"

	"PaperTrail-fm.com/db"
)

type Timeline struct {
	Id           string    `json:"id"`
	RootPapperID string    `json:"root_papper_id"`
	Date         time.Time `json:"date"`
}

func (t *Timeline) Save() error {
	var timelineID string

	// Verifica se a timeline já existe no banco de dados
	query := `SELECT id FROM timelines WHERE id = $1`
	err := db.DB.QueryRow(query, t.Id).Scan(&timelineID)

	if err != nil && err != sql.ErrNoRows {
		// Se ocorrer um erro diferente de "sem linhas encontradas", retorna o erro
		return fmt.Errorf("error checking timeline existence: %v", err)
	}

	if err == sql.ErrNoRows {
		// Se não há linhas (timeline não existe), insere um novo registro
		insertQuery := `
		INSERT INTO timelines(id, root_papper_id, date) 
		VALUES ($1, $2, $3)`
		_, err := db.DB.Exec(insertQuery, t.Id, t.RootPapperID, t.Date)
		if err != nil {
			return fmt.Errorf("error inserting timeline: %v", err)
		}

		fmt.Println("Inserted new timeline into database")
	} else {
		fmt.Println("Timeline already exists in database")
	}

	return nil
}

func GetTimelineByID(id string) (*Timeline, error) {
	query := "SELECT id, root_papper_id, date FROM timelines WHERE id = $1"
	row := db.DB.QueryRow(query, id)

	var timeline Timeline
	err := row.Scan(&timeline.Id, &timeline.RootPapperID, &timeline.Date)
	if err != nil {
		return nil, err
	}

	return &timeline, nil
}

func (t *Timeline) Delete() error {
	query := `SELECT id FROM timelines WHERE id = $1`
	err := db.DB.QueryRow(query, t.Id).Scan(&t.Id)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking timeline existence: %v", err)
	}

	query = `DELETE FROM timelines WHERE id = $1`
	_, err = db.DB.Exec(query, t.Id)
	if err != nil {
		return fmt.Errorf("error removing timeline: %v", err)
	}

	fmt.Println("Removed timeline from database")
	return nil
}

func (t *Timeline) AddChapter(chapterID string) error {
	query := `UPDATE chapters SET timeline_id = $1 WHERE id = $2`
	_, err := db.DB.Exec(query, t.Id, chapterID)
	if err != nil {
		return fmt.Errorf("error adding chapter to timeline: %v", err)
	}

	fmt.Println("Added chapter to timeline")
	return nil
}

func (t *Timeline) RemoveChapter(chapterID string) error {
	query := `UPDATE chapters SET timeline_id = NULL WHERE id = $1`
	_, err := db.DB.Exec(query, chapterID)
	if err != nil {
		return fmt.Errorf("error removing chapter from timeline: %v", err)
	}

	fmt.Println("Removed chapter from timeline")
	return nil
}

func (t *Timeline) AddEvent(eventID string) error {
	query := `UPDATE events SET timeline_id = $1 WHERE id = $2`
	_, err := db.DB.Exec(query, t.Id, eventID)
	if err != nil {
		return fmt.Errorf("error adding event to timeline: %v", err)
	}

	fmt.Println("Added event to timeline")
	return nil
}

func (t *Timeline) RemoveEvent(eventID string) error {
	query := `UPDATE events SET timeline_id = NULL WHERE id = $1`
	_, err := db.DB.Exec(query, eventID)
	if err != nil {
		return fmt.Errorf("error removing event from timeline: %v", err)
	}

	fmt.Println("Removed event from timeline")
	return nil
}
