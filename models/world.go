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
	chapters := []Chapter{}
	events := []Event{}
	connections := []Connection{}
	timelines := []Timeline{}
	storyLines := []StoryLine{}
	pappers := []Papper{}
	env := Env{}

	query := `SELECT name, created_at FROM worlds WHERE id = $1`
	err := db.DB.QueryRow(query, rp.Id).Scan(&rp.Name, &rp.CreatedAt)

	if err != nil && err != sql.ErrNoRows {
		return env, fmt.Errorf("error getting world info: %v", err)
	}

	eventsQuery := `
		SELECT id, name, start_date, end_date, world_id
		FROM events
		WHERE world_id = $1
	`
	rows, err := db.DB.Query(eventsQuery, rp.Id)
	if err != nil {
		return env, err
	}
	defer rows.Close()

	for rows.Next() {
		var event Event
		if err := rows.Scan(&event.Id, &event.Name, &event.Start_date, &event.End_date, &event.World_id); err != nil {
			return env, err
		}
		events = append(events, event)
	}

	// Consultar capítulos junto com as informações da tabela chapter_timeline
	chaptersQuery := `
    SELECT c.id, c.name, c.description, 
        c.created_at, c.papper_id, c.world_id, 
        c.event_id, c.storyline_id, c.timeline_id, c.order,
        ct.range
    FROM chapters c
    LEFT JOIN chapter_timeline ct ON c.id = ct.chapter_id
    WHERE c.world_id = $1
`
	rows, err = db.DB.Query(chaptersQuery, rp.Id)
	if err != nil {
		return env, err
	}
	defer rows.Close()

	for rows.Next() {
		var chapter Chapter
		var chapterRange sql.NullInt32 // Para lidar com valores nulos de range da tabela chapter_timeline

		// Faz o Scan para capturar os campos da tabela chapters e chapter_timeline
		if err := rows.Scan(&chapter.Id, &chapter.Name, &chapter.Description, &chapter.CreatedAt,
			&chapter.PapperID, &chapter.WorldsID, &chapter.EventID, &chapter.Storyline_id, &chapter.TimelineID, &chapter.Order, &chapterRange); err != nil {
			return env, err
		}

		// Se houver um valor válido para o range, atribuímos ao capítulo
		if chapterRange.Valid {
			chapter.Range = int(chapterRange.Int32)
		}

		chapters = append(chapters, chapter)
	}

	// Consultar conexões
	connectionsQuery := `
		SELECT id, source_chapter_id, target_chapter_id
		FROM connections
		WHERE source_chapter_id IN (SELECT id FROM chapters WHERE world_id = $1)
	`
	rows, err = db.DB.Query(connectionsQuery, rp.Id)
	if err != nil {
		return env, err
	}
	defer rows.Close()

	for rows.Next() {
		var connection Connection
		if err := rows.Scan(&connection.Id, &connection.SourceChapterID, &connection.TargetChapterID); err != nil {
			return env, err
		}
		connection.World_id = rp.Id
		connections = append(connections, connection)
	}

	// Consultar timelines
	timelinesQuery := `
		SELECT id, name, world_id, range, "order", description
		FROM timelines
		WHERE world_id = $1
	`
	rows, err = db.DB.Query(timelinesQuery, rp.Id)
	if err != nil {
		return env, err
	}
	defer rows.Close()

	for rows.Next() {
		var timeline Timeline
		if err := rows.Scan(&timeline.Id, &timeline.Name, &timeline.WorldsID, &timeline.Range, &timeline.Order, &timeline.Description); err != nil {
			return env, err
		}
		timelines = append(timelines, timeline)
	}
	// Consultar storyLines
	storylinesQuery := `
		SELECT id, name, "order", description
		FROM storylines
		WHERE world_id = $1
	`
	rows, err = db.DB.Query(storylinesQuery, rp.Id)
	if err != nil {
		return env, err
	}
	defer rows.Close()

	for rows.Next() {
		var storyLine StoryLine
		if err := rows.Scan(&storyLine.Id, &storyLine.Name, &storyLine.Order, &storyLine.Description); err != nil {
			return env, err
		}
		storyLines = append(storyLines, storyLine)
	}
	// Consultar timelines
	papperQuery := `
		SELECT id, name, description, created_at, "order"
		FROM pappers
		WHERE world_id = $1
	`
	rows, err = db.DB.Query(papperQuery, rp.Id)
	if err != nil {
		return env, err
	}
	defer rows.Close()

	for rows.Next() {
		var papper Papper
		if err := rows.Scan(&papper.ID, &papper.Name, &papper.Description, &papper.Created_at, &papper.Order); err != nil {
			return env, err
		}

		pappers = append(pappers, papper)
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
