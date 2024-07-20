package models

import (
	"database/sql"
	"fmt"
	"time"

	"PaperTrail-fm.com/db"
)

type Chapter struct {
	Id             string    `json:"id"`
	Root_papper_id string    `json:"root_papper_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Created_at     time.Time `json:"created_at"`
	Papper_id      string    `json:"papper_id"`
}

func (e *Chapter) Save() error {
	var papperId int

	// Verifica se o papper já existe no banco de dados
	query := `SELECT id FROM chapters WHERE id = $1 AND Papper_id = $2`
	err := db.DB.QueryRow(query, e.Id, e.Papper_id).Scan(&papperId)

	if err != nil && err != sql.ErrNoRows {
		// Se ocorrer um erro diferente de "sem linhas encontradas", retorna o erro
		return fmt.Errorf("error checking papper existence: %v", err)
	}

	if err == sql.ErrNoRows {
		// Se não há linhas (papper não existe), insere um novo registro
		insertQuery := `
		INSERT INTO chapters(name, description, Papper_id, id, created_at) 
		VALUES ($1, $2, $3, $4, $5)`
		_, err := db.DB.Exec(insertQuery, e.Name, e.Description, e.Papper_id, e.Id, e.Created_at)
		if err != nil {
			return fmt.Errorf("error inserting papper: %v", err)
		}

		fmt.Println("Inserted new papper into database")
	} else {
		fmt.Println("Papper already exists in database")
	}

	return nil
}

func GetChapterByID(id string) (*Chapter, error) {
	query := "SELECT * FROM chapters WHERE id = ?"
	row := db.DB.QueryRow(query, id)

	var chapter Chapter
	err := row.Scan(&chapter.Id, &chapter.Name, &chapter.Description, &chapter.Created_at, &chapter.Papper_id)
	if err != nil {
		return nil, err
	}

	return &chapter, nil
}
