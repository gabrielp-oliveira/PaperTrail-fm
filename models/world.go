package models

import (
	"database/sql"
	"fmt"
	"time"

	"PaperTrail-fm.com/db"
)

type World struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	UserID      string    `json:"user_id"`
	Description string    `json:"description"`
}

type Env struct {
	World
	Chapters         []Chapter         `json:"chapters"`
	Events           []Event           `json:"events"`
	Connections      []Connection      `json:"connections"`
	Timelines        []Timeline        `json:"timelines"`
	Papers           []Paper           `json:"papers"`
	StoryLines       []StoryLine       `json:"storyLines"`
	GroupConnections []GroupConnection `json:"groupConnections"`
	Subway_settings  Subway_settings   `json:"subway_settings"`
}

// Salva o Worlds no banco de dados
func (rp *World) Save() error {
	var rpId string

	query := `SELECT id FROM worlds WHERE id = $1`
	err := db.DB.QueryRow(query, rp.Id).Scan(&rpId)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking worlds existence: %v", err)
	}

	if err == sql.ErrNoRows {
		insertQuery := `
		INSERT INTO worlds(id, name, description, created_at, user_id) 
		VALUES ($1, $2, $3, $4, $5)`
		_, err := db.DB.Exec(insertQuery, rp.Id, rp.Name, rp.Description, rp.CreatedAt, rp.UserID)
		if err != nil {
			return fmt.Errorf("error inserting world: %v", err)
		}

		fmt.Println("Inserted new world into database")
	} else {
		fmt.Println("World already exists in database")
	}

	return nil
}
func (w *World) Get() error {

	query := "SELECT id, name, description, created_at FROM worlds WHERE id = $1"

	row := db.DB.QueryRow(query, w.Id)

	err := row.Scan(&w.Id, &w.Name, &w.Description, &w.CreatedAt)
	if err != nil {
		return err
	}

	return nil

}

// Obtém a lista de Papers associados ao Worlds
func (rp *World) GetPaperList() ([]Paper, error) {
	query := "SELECT id, name, description, created_at FROM Papers WHERE world_id = $1"
	rows, err := db.DB.Query(query, rp.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Paper
	for rows.Next() {
		var paper Paper
		if err := rows.Scan(&paper.ID, &paper.Name, &paper.Description, &paper.Created_at); err != nil {
			return nil, err
		}
		list = append(list, paper)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

// Obtém capítulos, conexões, eventos e timelines associados ao Worlds
func (rp *World) GetWorldData() (Env, error) {

	env := Env{}

	chapters, err := GetChapteJoinTimelineByWorldId(rp.Id)
	if err != nil {
		return env, err
	}
	timelines, err := GetTimelinesByWorldId(rp.Id)
	if err != nil {
		return env, err
	}

	storyLines, err := GetstoryLinesByWorldId(rp.Id)
	if err != nil {
		return env, err
	}

	Papers, err := GetPaperListByWorldId(rp.Id)
	if err != nil {
		return env, err
	}
	connections, err := GetConnectionsListByWorldId(rp.Id)
	if err != nil {
		return env, err
	}
	events, err := getEventsByWorldId(rp.Id)
	if err != nil {
		return env, err
	}
	ss, err := getSubwaySettingssByWorldId(rp.Id)
	if err != nil {
		return env, err
	}
	gcs, err := GetConnectionsGroupListByWorldId(rp.Id)
	if err != nil {
		return env, err
	}

	// Preencher o objeto Env
	env.World = *rp
	env.Chapters = chapters
	env.Connections = connections
	env.Events = events
	env.Timelines = timelines
	env.Papers = Papers
	env.StoryLines = storyLines
	env.Subway_settings = *ss
	env.GroupConnections = gcs
	return env, nil
}

func (rp *World) GetWorldChapters() ([]Chapter, error) {

	chapters, err := GetChapteJoinTimelineByWorldId(rp.Id)
	if err != nil {
		return chapters, err
	}

	return chapters, nil
}
