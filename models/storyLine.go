package models

import (
	"database/sql"
	"fmt"
	"time"

	"PaperTrail-fm.com/db"
	"github.com/google/uuid"
)

type StoryLine struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created_at  time.Time `json:"created_at"`
	WorldsID    string    `json:"world_id"`
	Order       int       `json:"order"`
}

func (t *StoryLine) Save() error {
	var storyLinesID string

	// Verifica se a storyLines já existe no banco de dados
	query := `SELECT id FROM storyLines WHERE id = $1`
	err := db.DB.QueryRow(query, t.Id).Scan(&storyLinesID)

	if err != nil && err != sql.ErrNoRows {
		// Se ocorrer um erro diferente de "sem linhas encontradas", retorna o erro
		return fmt.Errorf("error checking storyLines existence: %v", err)
	}

	if err == sql.ErrNoRows {
		var maxOrder *int
		orderQuery := `SELECT MAX("order") FROM storyLines WHERE world_id = $1 `
		err = db.DB.QueryRow(orderQuery, t.WorldsID).Scan(&maxOrder)
		if err != nil {
			return fmt.Errorf("error getting max order: %v", err)
		}

		newOrder := 1
		if maxOrder != nil {
			newOrder = *maxOrder + 1
		} else {
			newOrder = 1
		}

		t.Order = newOrder

		// Se não há linhas (storyLines não existe), insere um novo registro
		insertQuery := `
		INSERT INTO storyLines(id, name, description, created_at, world_id, "order") 
		VALUES ($1, $2, $3, $4, $5, $6)`
		_, err := db.DB.Exec(insertQuery, t.Id, t.Name, t.Description, t.Created_at, t.WorldsID, t.Order)
		if err != nil {
			return fmt.Errorf("error inserting storyLines: %v", err)
		}

		fmt.Println("Inserted new storyLines into database")
	} else {
		fmt.Println("storyLines already exists in database")
	}

	return nil
}

func GetstoryLinesByWorldId(worldId string) ([]StoryLine, error) {
	storyLines := []StoryLine{}

	storylinesQuery := `
		SELECT id, name, "order", description
		FROM storylines
		WHERE world_id = $1
	`
	rows, err := db.DB.Query(storylinesQuery, worldId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var storyLine StoryLine
		if err := rows.Scan(&storyLine.Id, &storyLine.Name, &storyLine.Order, &storyLine.Description); err != nil {
			return nil, err
		}
		storyLines = append(storyLines, storyLine)
	}

	return storyLines, err
}

func (t *StoryLine) Delete() error {
	query := `SELECT id FROM storyLines WHERE id = $1`
	err := db.DB.QueryRow(query, t.Id).Scan(&t.Id)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking storyLines existence: %v", err)
	}

	query = `DELETE FROM storyLines WHERE id = $1`
	_, err = db.DB.Exec(query, t.Id)
	if err != nil {
		return fmt.Errorf("error removing storyLines: %v", err)
	}

	fmt.Println("Removed storyLines from database")
	return nil
}

func (t StoryLine) CreateBasicStoryLines(wiD string) ([]StoryLine, error) {
	var storyLinesList []StoryLine

	var str1 StoryLine
	str1.Name = "Main"
	str1.Description = "Main story line"
	str1.WorldsID = wiD
	str1.Order = 1
	id := uuid.New().String()
	str1.Id = id
	str1.Created_at = time.Now()
	storyLinesList = append(storyLinesList, str1)
	err := str1.Save()
	if err != nil {
		return nil, fmt.Errorf("error removing storyLines: %v", err)
	}

	var str2 StoryLine
	str2.Name = "Seccondary"
	str2.Description = "Seccondary story line"
	str2.WorldsID = wiD
	str2.Order = 2
	id = uuid.New().String()
	str2.Created_at = time.Now()
	str2.Id = id
	err = str2.Save()
	if err != nil {
		return nil, fmt.Errorf("error removing storyLines: %v", err)
	}
	storyLinesList = append(storyLinesList, str2)
	var str3 StoryLine
	str3.Name = "extra"
	str3.Description = "extra story line"
	str3.WorldsID = wiD
	str3.Order = 3
	id = uuid.New().String()
	str3.Id = id
	str3.Created_at = time.Now()
	err = str3.Save()
	if err != nil {
		return nil, fmt.Errorf("error removing storyLines: %v", err)
	}
	storyLinesList = append(storyLinesList, str3)
	return storyLinesList, nil

}

func (t *StoryLine) Update() error {
	// id, name, description, range, "order"

	query := `
	UPDATE storylines
	SET name = $1, description = $2, "order" = $3
	WHERE id = $4
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(t.Name, t.Description, t.Order, t.Id)
	return err
}
