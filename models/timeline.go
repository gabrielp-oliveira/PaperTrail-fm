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

		var maxOrder *int
		orderQuery := `SELECT MAX("order") FROM timelines WHERE world_id = $1 `
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

func GetTimelinesByWorldId(worldId string) ([]Timeline, error) {
	timelines := []Timeline{}
	timelinesQuery := `
		SELECT id, name, world_id, range, "order", description
		FROM timelines
		WHERE world_id = $1 ORDER BY "order"
	`
	rows, err := db.DB.Query(timelinesQuery, worldId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var timeline Timeline
		if err := rows.Scan(&timeline.Id, &timeline.Name, &timeline.WorldsID, &timeline.Range, &timeline.Order, &timeline.Description); err != nil {
			return nil, err
		}
		timelines = append(timelines, timeline)
	}
	return timelines, err
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

func (tl Timeline) Create() error {
	timelineID := uuid.New().String()
	var desc Description
	desc.Description_data = tl.Description

	tl.Id = timelineID
	descId := uuid.New().String()
	tl.Description = descId

	err := tl.Save()
	desc.Id = descId
	desc.Resource_type = "timeline"
	desc.Resource_id = timelineID
	if err != nil {
		return err
	}
	err = desc.Save()
	if err != nil {
		return err
	}
	return nil
}
func (t Timeline) CreateBasicTimelines(wiD string) ([]Timeline, error) {
	var TimelineList []Timeline

	var tmLine1 Timeline
	tmLine1.Name = "start"
	tmLine1.Description = "start time line"
	tmLine1.WorldsID = wiD
	tmLine1.Order = 1
	tmLine1.Range = 15
	tmLine1.Id = uuid.New().String()
	err := tmLine1.Create()
	if err != nil {
		return nil, fmt.Errorf("error creating Timeline: %v", err)
	}
	TimelineList = append(TimelineList, tmLine1)

	var tmLine2 Timeline
	tmLine2.Name = "seccond"
	tmLine2.Description = "Seccondary time line"
	tmLine2.WorldsID = wiD
	tmLine2.Order = 2
	tmLine2.Range = 15
	tmLine2.Id = uuid.New().String()
	err = tmLine2.Create()
	if err != nil {
		return nil, fmt.Errorf("error creating Timeline: %v", err)
	}

	TimelineList = append(TimelineList, tmLine2)
	var tmLine3 Timeline
	tmLine3.Name = "third"
	tmLine3.Description = "third time line"
	tmLine3.WorldsID = wiD
	tmLine3.Order = 3
	tmLine3.Range = 15

	tmLine3.Id = uuid.New().String()
	err = tmLine3.Create()
	if err != nil {
		return nil, fmt.Errorf("error removing Timeline: %v", err)
	}
	TimelineList = append(TimelineList, tmLine2)
	var tmLine4 Timeline
	tmLine4.Name = "fourth"
	tmLine4.Description = "fourth time line"
	tmLine4.WorldsID = wiD
	tmLine4.Order = 4
	tmLine4.Range = 15

	tmLine4.Id = uuid.New().String()
	err = tmLine4.Create()
	if err != nil {
		return nil, fmt.Errorf("error removing Timeline: %v", err)
	}
	TimelineList = append(TimelineList, tmLine4)
	var tmLine5 Timeline
	tmLine5.Name = "fifth"
	tmLine5.Description = "fifth time line"
	tmLine5.WorldsID = wiD
	tmLine5.Order = 5
	tmLine5.Range = 15

	tmLine5.Id = uuid.New().String()
	err = tmLine5.Create()
	if err != nil {
		return nil, fmt.Errorf("error removing Timeline: %v", err)
	}
	TimelineList = append(TimelineList, tmLine5)
	return TimelineList, nil

}

func (t *Timeline) Update() error {
	// id, name, description, range, "order"

	query := `
	UPDATE timelines
	SET name = $1, description = $2,range= $3, "order" = $4	WHERE id = $5
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(t.Name, t.Description, t.Range, t.Order, t.Id)
	return err
}

func (t *Timeline) Get() error {
	query := `
		SELECT id, name, world_id, range, "order", description
		FROM timelines
		WHERE id = $1
	`
	row := db.DB.QueryRow(query, t.Id)

	err := row.Scan(&t.Id, &t.Name, &t.WorldsID, &t.Range, &t.Order, &t.Description)

	if err != nil {
		return err
	}
	return nil
}
