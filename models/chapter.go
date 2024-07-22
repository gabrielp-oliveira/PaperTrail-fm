package models

import (
	"database/sql"
	"fmt"
	"time"

	"PaperTrail-fm.com/db"
)

type Chapter struct {
	Id           string         `json:"id"`
	RootPapperID string         `json:"root_papper_id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	CreatedAt    time.Time      `json:"created_at"`
	PapperID     string         `json:"papper_id"`
	EventID      sql.NullString `json:"event_id"`
	TimelineID   sql.NullString `json:"timeline_id"`
}

func (c *Chapter) Save() error {
	var chapterID string

	// Verifica se o capítulo já existe no banco de dados
	query := `SELECT id FROM chapters WHERE id = $1`
	err := db.DB.QueryRow(query, c.Id).Scan(&chapterID)

	if err != nil && err != sql.ErrNoRows {
		// Se ocorrer um erro diferente de "sem linhas encontradas", retorna o erro
		return fmt.Errorf("error checking chapter existence: %v", err)
	}

	if err == sql.ErrNoRows {
		// Se não há linhas (capítulo não existe), insere um novo registro
		insertQuery := `
		INSERT INTO chapters(id, name, description, created_at, papper_id, root_papper_id, event_id, timeline_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		_, err := db.DB.Exec(insertQuery, c.Id, c.Name, c.Description, c.CreatedAt, c.PapperID, c.RootPapperID, c.EventID, c.TimelineID)
		if err != nil {
			return fmt.Errorf("error inserting chapter: %v", err)
		}

		fmt.Println("Inserted new chapter into database")
	} else {
		fmt.Println("Chapter already exists in database")
	}

	return nil
}

func GetChapterByID(id string) (*Chapter, error) {
	query := "SELECT id, name, description, created_at, papper_id, root_papper_id, event_id, timeline_id FROM chapters WHERE id = $1"
	row := db.DB.QueryRow(query, id)

	var chapter Chapter
	err := row.Scan(&chapter.Id, &chapter.Name, &chapter.Description, &chapter.CreatedAt, &chapter.PapperID, &chapter.RootPapperID, &chapter.EventID, &chapter.TimelineID)
	if err != nil {
		return nil, err
	}

	return &chapter, nil
}

func (c *Chapter) AddEvent(eventID string) error {
	query := `UPDATE chapters SET event_id = $1 WHERE id = $2`
	_, err := db.DB.Exec(query, eventID, c.Id)
	if err != nil {
		return fmt.Errorf("error adding event to chapter: %v", err)
	}

	fmt.Println("Added event to chapter")
	return nil
}

func (c *Chapter) RemoveEvent() error {
	query := `UPDATE chapters SET event_id = NULL WHERE id = $1`
	_, err := db.DB.Exec(query, c.Id)
	if err != nil {
		return fmt.Errorf("error removing event from chapter: %v", err)
	}

	fmt.Println("Removed event from chapter")
	return nil
}

func (c *Chapter) AddTimeline(timelineID string) error {
	query := `UPDATE chapters SET timeline_id = $1 WHERE id = $2`
	_, err := db.DB.Exec(query, timelineID, c.Id)
	if err != nil {
		return fmt.Errorf("error adding timeline to chapter: %v", err)
	}

	fmt.Println("Added timeline to chapter")
	return nil
}

func (c *Chapter) RemoveTimeline() error {
	query := `UPDATE chapters SET timeline_id = NULL WHERE id = $1`
	_, err := db.DB.Exec(query, c.Id)
	if err != nil {
		return fmt.Errorf("error removing timeline from chapter: %v", err)
	}

	fmt.Println("Removed timeline from chapter")
	return nil
}
