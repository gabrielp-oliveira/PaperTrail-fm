package models

import (
	"database/sql"
	"fmt"
	"time"

	"PaperTrail-fm.com/db"
)

type ChapterDetails struct {
	Chapter         `json:"chapter"`
	Timeline        `json:"timeline"`
	StoryLine       `json:"storyLine"`
	DocumentUrl     string              `json:"documentUrl"`
	Color           string              `json:"color"`
	Events          []Event             `json:"events"`
	Next            Chapter             `json:"next"`
	Prev            Chapter             `json:"prev"`
	RelatedChapters []ChapterConnection `json:"relatedChapters"`
}
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

type RelatedChapters struct {
	ChapterID      string
	RelatedChapter string
	GroupName      string
	GroupColor     string
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

func (c *Chapter) Get() (string, error) {
	var color string
	query := `
		SELECT 
			ch.id, 
			ch.name, 
			ch.description, 
			ch.created_at, 
			ch.Paper_id, 
			ch.world_id, 
			ch.event_id, 
			ch.timeline_id, 
			ch.update, 
			ch.range, 
			ch."order", 
			ch.last_update, 
			ch.storyline_id,
			p.color
		FROM chapters ch
		LEFT JOIN Papers p ON ch.Paper_id = p.id
		WHERE ch.id = $1
	`

	row := db.DB.QueryRow(query, c.Id)

	// Adicione um campo "Color" no struct do Chapter, por exemplo: Color string
	err := row.Scan(&c.Id, &c.Name, &c.Description, &c.CreatedAt, &c.PaperID, &c.WorldsID, &c.Event_Id, &c.TimelineID, &c.Update, &c.Range, &c.Order, &c.LastUpdate, &c.Storyline_id, &color)
	if err != nil {
		return color, err
	}
	return color, nil
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

func GetChapteJoinTimelineByWorldId(worldId string) ([]Chapter, error) {
	chapters := []Chapter{}

	chaptersQuery := `
    SELECT id, name, description, 
        created_at, Paper_id, world_id, 
        event_id, storyline_id, timeline_id, "order",
        range
    FROM chapters
    WHERE world_id = $1
`
	rows, err := db.DB.Query(chaptersQuery, worldId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var chapter Chapter

		// Faz o Scan para capturar os campos da tabela chapters
		if err := rows.Scan(&chapter.Id, &chapter.Name, &chapter.Description, &chapter.CreatedAt,
			&chapter.PaperID, &chapter.WorldsID, &chapter.Event_Id, &chapter.Storyline_id, &chapter.TimelineID, &chapter.Order, &chapter.Range); err != nil {
			return nil, err
		}

		chapters = append(chapters, chapter)
	}
	return chapters, nil
}

func (c *Chapter) GetChapterByPaperAndOrder(order int) (*Chapter, error) {
	query := `SELECT 
		id, 
		name, 
		description, 
		created_at, 
		Paper_id, 
		world_id, 
		event_id, 
		timeline_id, 
		update, 
		"order", 
		last_update, 
		storyline_id
	 	FROM chapters WHERE "order" = $1 and paper_id = $2`

	row := db.DB.QueryRow(query, order, c.PaperID)

	var chapter Chapter
	err := row.Scan(&chapter.Id, &chapter.Name, &chapter.Description, &chapter.CreatedAt, &chapter.PaperID, &chapter.WorldsID, &chapter.Event_Id, &chapter.TimelineID, &chapter.Update, &chapter.Order, &chapter.LastUpdate, &chapter.Storyline_id)
	if err != nil {
		return nil, err
	}

	return &chapter, nil
}

type ChapterConnection struct {
	ChapterID      string `json:"chapterId"`
	RelatedChapter struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"relatedChapter"`
	GroupName  string `json:"groupName"`
	GroupColor string `json:"groupColor"`
}

func (c *Chapter) GetRelatedChapters() ([]ChapterConnection, error) {
	query := `
		SELECT
			CASE
				WHEN c.source_chapter_id = $1 THEN ch_target.id
				WHEN c.target_chapter_id = $1 THEN ch_source.id
			END AS related_chapter_id,
			CASE
				WHEN c.source_chapter_id = $1 THEN ch_target.name
				WHEN c.target_chapter_id = $1 THEN ch_source.name
			END AS related_chapter_name,
			COALESCE(g.name, '') AS group_name,
			COALESCE(g.color, '#000000') AS group_color
		FROM
			connections c
		LEFT JOIN
			chapters ch_source ON c.source_chapter_id = ch_source.id
		LEFT JOIN
			chapters ch_target ON c.target_chapter_id = ch_target.id
		LEFT JOIN
			group_connections g ON c.group_id = g.id
		WHERE
			c.source_chapter_id = $1 OR c.target_chapter_id = $1;
	`

	rows, err := db.DB.Query(query, c.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var connections []ChapterConnection
	for rows.Next() {
		var connection ChapterConnection
		var relatedChapterID, relatedChapterName string

		connection.ChapterID = c.Id
		if err := rows.Scan(
			&relatedChapterID,
			&relatedChapterName,
			&connection.GroupName,
			&connection.GroupColor,
		); err != nil {
			return nil, err
		}

		connection.RelatedChapter.ID = relatedChapterID
		connection.RelatedChapter.Name = relatedChapterName

		connections = append(connections, connection)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return connections, nil
}
