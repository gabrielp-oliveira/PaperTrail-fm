package models

import (
	"database/sql"
	"fmt"
	"time"

	"PaperTrail-fm.com/db"
)

type World struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UserID    string    `json:"user_id"`
}

type Env struct {
	World
	Chapters    []Chapter    `json:"chapters"`
	Events      []Event      `json:"events"`
	Connections []Connection `json:"connections"`
	Timelines   []Timeline   `json:"timelines"`
	Pappers     []Papper     `json:"pappers"`
	StoryLines  []StoryLine  `json:"storyLines"`
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
func (rp *World) GetPapperList() ([]Papper, error) {
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

	pappers, err := GetPapperListByWorldId(rp.Id)
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

	// Preencher o objeto Env
	env.World = *rp
	env.Chapters = chapters
	env.Connections = connections
	env.Events = events
	env.Timelines = timelines
	env.Pappers = pappers
	env.StoryLines = storyLines
	return env, nil
}
