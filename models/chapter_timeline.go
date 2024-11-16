package models

import (
	"database/sql"
	"fmt"

	"PaperTrail-fm.com/db"
	"github.com/google/uuid"
)

type ChapterTimeline struct {
	Id         string `json:"id"`
	Chapter_Id string `json:"chapter_Id"`
	TimelineID string `json:"timeline_id"`
	Range      int    `json:"range"`
}

type ChapterTl struct {
	Chapter
	Range        int    `json:"range"`
	Chapter_Id   string `json:"chapter_Id"`
	TimelineID   string `json:"timeline_id"`
	Storyline_id string `json:"storyline_id"`
	Order        int    `json:"order"`
	Event_Id     string `json:"event_id"`
}

func (t *ChapterTimeline) Save() error {
	var ranges int

	// Verifica se a Timeline já existe no banco de dados
	query := `SELECT range FROM chapter_timeline WHERE chapter_Id = $1 and TimeLine_Id = $2`
	err := db.DB.QueryRow(query, t.Id).Scan(&ranges)

	if err != nil && err != sql.ErrNoRows {
		// Se ocorrer um erro diferente de "sem linhas encontradas", retorna o erro
		return fmt.Errorf("error checking chapter_timeline existence: %v", err)
	}

	if err == sql.ErrNoRows {
		// Se não há linhas (Timeline não existe), insere um novo registro
		insertQuery := `
		INSERT INTO chapter_timeline(id, chapter_Id, range, TimeLine_Id) 
		VALUES ($1, $2, $3, $4)`
		_, err := db.DB.Exec(insertQuery, t.Id, t.Chapter_Id, t.Range, t.TimelineID)
		if err != nil {
			return fmt.Errorf("error inserting Timeline: %v", err)
		}

		fmt.Println("Inserted new Timeline into database")
	} else {
		fmt.Println("Timeline already exists in database")
	}

	return nil
}

func (t *ChapterTimeline) Delete() error {
	query := `SELECT id FROM timelines WHERE id = $1`
	err := db.DB.QueryRow(query, t.Id).Scan(&t.Id)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking Timeline existence: %v", err)
	}

	query = `DELETE FROM timelines WHERE id = $1`
	_, err = db.DB.Exec(query, t.Id)
	if err != nil {
		return fmt.Errorf("error removing Timeline: %v", err)
	}

	fmt.Println("Removed Timeline from database")
	return nil
}

func (t *ChapterTimeline) GetRangeByChapterId() error {
	query := `SELECT range FROM chapter_timeline WHERE chapter_Id = $1`
	err := db.DB.QueryRow(query, t.Id).Scan(&t.Range)

	if err != nil && err != sql.ErrNoRows {
		// Se ocorrer um erro diferente de "sem linhas encontradas", retorna o erro
		return fmt.Errorf("error checking chapter_timeline existence: %v", err)
	}
	return nil
}

func (t *ChapterTimeline) Update() error {

	var id string

	// Verifica se a Timeline já existe no banco de dados
	query := `SELECT id FROM chapter_timeline WHERE chapter_Id = $1`
	err := db.DB.QueryRow(query, t.Chapter_Id).Scan(&id)

	if err != nil && err != sql.ErrNoRows {
		// Se ocorrer um erro diferente de "sem linhas encontradas", retorna o erro
		return fmt.Errorf("error checking chapter_timeline existence: %v", err)
	}

	if err == sql.ErrNoRows {
		// Se não há linhas (Timeline não existe), insere um novo registro
		id = uuid.New().String()

		insertQuery := `
		INSERT INTO chapter_timeline(id, chapter_Id, range, TimeLine_Id) 
		VALUES ($1, $2, $3, $4)`
		_, err := db.DB.Exec(insertQuery, id, t.Chapter_Id, t.Range, t.TimelineID)
		if err != nil {
			return fmt.Errorf("error inserting Timeline: %v", err)
		}

		fmt.Println("Inserted new Timeline into database")
	} else {

		query := `
		UPDATE chapter_timeline
		SET "range" = $1, TimeLine_Id = $2
		WHERE chapter_Id = $3
		`
		stmt, err := db.DB.Prepare(query)

		if err != nil {
			return err
		}

		defer stmt.Close()

		_, err = stmt.Exec(t.Range, t.TimelineID, t.Chapter_Id)
		return err
	}

	return nil

}
