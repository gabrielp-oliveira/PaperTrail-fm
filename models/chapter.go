package models

import (
	"database/sql"
	"fmt"
	"time"

	"PaperTrail-fm.com/db"
)

type Chapter struct {
	Id           string     `json:"id"`
	WorldsID     string     `json:"world_id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	CreatedAt    time.Time  `json:"created_at"`
	PaperID      string     `json:"paper_id"`
	Event_Id     *string    `json:"event_Id"`
	TimelineID   *string    `json:"timeline_id"`
	Storyline_id *string    `json:"storyline_id"`
	Link         string     `json:"link"`
	Update       *string    `json:"update"`
	Order        int        `json:"order"`
	Range        int        `json:"range"`
	LastUpdate   *time.Time `json:"last_update"` // Usa ponteiro para time.Time
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

		var maxOrder *int
		orderQuery := `SELECT MAX("order") FROM chapters WHERE world_id = $1 and Paper_id = $2`
		err := db.DB.QueryRow(orderQuery, c.WorldsID, c.PaperID).Scan(&maxOrder)
		if err != nil {
			return fmt.Errorf("error getting max order: %v", err)
		}

		newOrder := 1
		if maxOrder != nil {
			newOrder = *maxOrder + 1
		} else {
			newOrder = 1

		}

		insertQuery := `
		INSERT INTO chapters(id, name, description, created_at, Paper_id, world_id, event_id, timeline_id, "order") 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
		_, err = db.DB.Exec(insertQuery, c.Id, c.Name, c.Description, c.CreatedAt, c.PaperID, c.WorldsID, c.Event_Id, c.TimelineID, newOrder)
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
	query := "SELECT id, name, description, created_at, Paper_id, world_id, event_id, timeline_id FROM chapters WHERE id = $1"
	row := db.DB.QueryRow(query, id)

	var chapter Chapter
	err := row.Scan(&chapter.Id, &chapter.Name, &chapter.Description, &chapter.CreatedAt, &chapter.PaperID, &chapter.WorldsID, &chapter.Event_Id, &chapter.TimelineID)
	if err != nil {
		return nil, err
	}

	return &chapter, nil
}

func (c *Chapter) AddEvent(Event_Id string) error {
	query := `UPDATE chapters SET event_id = $1 WHERE id = $2`
	_, err := db.DB.Exec(query, Event_Id, c.Id)
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

func (c *Chapter) Get() error {
	query := `SELECT id, name, description, created_at, Paper_id, world_id, event_id, timeline_id, update, "order", last_update FROM chapters WHERE id = $1`
	row := db.DB.QueryRow(query, c.Id)

	err := row.Scan(&c.Id, &c.Name, &c.Description, &c.CreatedAt, &c.PaperID, &c.WorldsID, &c.Event_Id, &c.TimelineID, &c.Update, &c.Order, &c.LastUpdate)
	if err != nil {
		return err
	}
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
func (c *Chapter) UpdateChapter() error {
	query := `
	UPDATE chapters
	SET name = $1, description = $2, "order" = $3, update = $4, last_update = $5, storyline_id = $6, timeline_id = $7, event_id = $8, range = $9
	WHERE id = $10
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	if c.TimelineID != nil && *c.TimelineID == "" {
		c.TimelineID = nil
	}
	if c.Storyline_id != nil && *c.Storyline_id == "" {
		c.Storyline_id = nil
	}
	if c.Event_Id != nil && *c.Event_Id == "" {
		c.Event_Id = nil
	}

	_, err = stmt.Exec(c.Name, c.Description, c.Order, c.Update, c.LastUpdate, c.Storyline_id, c.TimelineID, c.Event_Id, c.Range, c.Id)
	return err
}
