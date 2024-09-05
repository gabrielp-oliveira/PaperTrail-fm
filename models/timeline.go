package models

import (
	"database/sql"
	"fmt"

	"PaperTrail-fm.com/db"
	"github.com/google/uuid"
)

type Timeline struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	WorldsID    string `json:"world_id"`
	Order       int    `json:"order"`
	Range       int    `json:"range"`
}

func (t *Timeline) Save() error {
	var timelineId string

	// Verifica se a Timeline já existe no banco de dados
	query := `SELECT id FROM timelines WHERE id = $1`
	err := db.DB.QueryRow(query, t.Id).Scan(&timelineId)

	if err != nil && err != sql.ErrNoRows {
		// Se ocorrer um erro diferente de "sem linhas encontradas", retorna o erro
		return fmt.Errorf("error checking Timelines existence: %v", err)
	}

	if err == sql.ErrNoRows {
		// Se não há linhas (Timeline não existe), insere um novo registro
		insertQuery := `
		INSERT INTO timelines(id, name, description, range, world_id, "order") 
		VALUES ($1, $2, $3, $4, $5, $6)`
		_, err := db.DB.Exec(insertQuery, t.Id, t.Name, t.Description, t.Range, t.WorldsID, t.Order)
		if err != nil {
			return fmt.Errorf("error inserting Timeline: %v", err)
		}

		fmt.Println("Inserted new Timeline into database")
	} else {
		fmt.Println("Timeline already exists in database")
	}

	return nil
}

func GetTimelinesByWorldId(worldId string) (*Timeline, error) {
	query := "SELECT id, name, description, range, worldsId, 'order' FROM Timelines WHERE id = $1"
	row := db.DB.QueryRow(query, worldId)

	var Timeline Timeline
	err := row.Scan(&Timeline.Id, &Timeline.WorldsID, &Timeline.Name, &Timeline.Description, &Timeline.Range, &Timeline.WorldsID, &Timeline.Order)
	if err != nil {
		return nil, err
	}

	return &Timeline, nil
}

func (t *Timeline) Delete() error {
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

func (t Timeline) CreateBasicTimelines(wiD string) ([]Timeline, error) {
	var TimelineList []Timeline

	var tmLine1 Timeline
	tmLine1.Name = "start"
	tmLine1.Description = "start time line"
	tmLine1.WorldsID = wiD
	tmLine1.Order = 1
	tmLine1.Range = 5

	id := uuid.New().String()

	tmLine1.Id = id
	TimelineList = append(TimelineList, tmLine1)
	err := tmLine1.Save()
	if err != nil {
		return nil, fmt.Errorf("error removing Timeline: %v", err)
	}

	var tmLine2 Timeline
	tmLine2.Name = "seccond"
	tmLine2.Description = "Seccondary time line"
	tmLine2.WorldsID = wiD
	tmLine2.Order = 2
	tmLine2.Range = 5
	id = uuid.New().String()
	tmLine2.Id = id
	err = tmLine2.Save()
	if err != nil {
		return nil, fmt.Errorf("error removing Timeline: %v", err)
	}
	TimelineList = append(TimelineList, tmLine2)
	var tmLine3 Timeline
	tmLine3.Name = "third"
	tmLine3.Description = "third time line"
	tmLine3.WorldsID = wiD
	tmLine3.Order = 3
	tmLine3.Range = 5
	id = uuid.New().String()

	tmLine3.Id = id
	err = tmLine3.Save()
	if err != nil {
		return nil, fmt.Errorf("error removing Timeline: %v", err)
	}
	TimelineList = append(TimelineList, tmLine2)
	var tmLine4 Timeline
	tmLine4.Name = "fourth"
	tmLine4.Description = "fourth time line"
	tmLine4.WorldsID = wiD
	tmLine4.Order = 4
	tmLine4.Range = 5

	id = uuid.New().String()

	tmLine4.Id = id
	err = tmLine4.Save()
	if err != nil {
		return nil, fmt.Errorf("error removing Timeline: %v", err)
	}
	TimelineList = append(TimelineList, tmLine4)
	var tmLine5 Timeline
	tmLine5.Name = "fifth"
	tmLine5.Description = "fifth time line"
	tmLine5.WorldsID = wiD
	tmLine5.Order = 5
	tmLine5.Range = 5

	id = uuid.New().String()

	tmLine5.Id = id
	err = tmLine5.Save()
	if err != nil {
		return nil, fmt.Errorf("error removing Timeline: %v", err)
	}
	TimelineList = append(TimelineList, tmLine5)
	return TimelineList, nil

}
