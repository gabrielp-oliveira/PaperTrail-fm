package models

import (
	"database/sql"
	"fmt"
	"math"

	"PaperTrail-fm.com/db"
)

type Subway_settings struct {
	Id                              string  `json:"id"`
	Chapter_names                   bool    `json:"chapter_names"`
	Display_table_chapters          bool    `json:"display_table_chapters"`
	Timeline_update_chapter         bool    `json:"timeline_update_chapter"`
	Storyline_update_chapter        bool    `json:"storyline_update_chapter"`
	Zomm                            float64 `json:"zoom"`
	X                               float64 `json:"x"`
	Y                               float64 `json:"y"`
	World_id                        string  `json:"world_id"`
	Group_connection_update_chapter bool    `json:"group_connection_update_chapter"`
}

func roundToTwoDecimalPlaces(value float64) float64 {
	return math.Round(value*100) / 100
}
func (ss *Subway_settings) Save() error {
	var id int

	query := `SELECT id FROM subway_settings WHERE world_id = $1`
	err := db.DB.QueryRow(query, ss.Id).Scan(&id)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking subway_settings existence: %v", err)
	}

	if err == sql.ErrNoRows {

		insertQuery := `
		INSERT INTO subway_settings(id, Chapter_names, zoom, x, y,display_table_chapters,storyline_update_chapter,timeline_update_chapter, World_id, group_connection_update_chapter) 
		VALUES ($1, $2, $3,  $4, $5, $6, $7, $8, $9, $10)`
		_, err := db.DB.Exec(insertQuery, ss.Id, ss.Chapter_names, roundToTwoDecimalPlaces(ss.Zomm), roundToTwoDecimalPlaces(ss.X),
			roundToTwoDecimalPlaces(ss.Y), ss.Display_table_chapters, ss.Storyline_update_chapter, ss.Timeline_update_chapter, ss.World_id, ss.Group_connection_update_chapter)
		if err != nil {
			return fmt.Errorf("error inserting settings: %v", err)
		}
		fmt.Println("Inserted new settings into database")
	} else {
		fmt.Println("settings already exists in database")
	}
	return nil
}

func (ss *Subway_settings) Update() error {
	query := `
	UPDATE subway_settings
	SET Chapter_names = $1, zoom = $2, x = $3, y = $4, display_table_chapters = $5, storyline_update_chapter = $6 ,timeline_update_chapter = $7, Group_connection_update_chapter = $8
	WHERE id = $9
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(ss.Chapter_names, roundToTwoDecimalPlaces(ss.Zomm), roundToTwoDecimalPlaces(ss.X),
		roundToTwoDecimalPlaces(ss.Y), ss.Display_table_chapters, ss.Storyline_update_chapter, ss.Timeline_update_chapter, ss.Group_connection_update_chapter, ss.Id)
	return err
}

func getSubwaySettingssByWorldId(worldId string) (*Subway_settings, error) {
	var ss Subway_settings

	ssquery := `
	SELECT id, zoom, x, y, Chapter_names, display_table_chapters, storyline_update_chapter, timeline_update_chapter, group_connection_update_chapter
	FROM subway_settings
	WHERE world_id = $1
`
	row := db.DB.QueryRow(ssquery, worldId)

	err := row.Scan(&ss.Id, &ss.Zomm, &ss.X, &ss.Y, &ss.Chapter_names, &ss.Display_table_chapters, &ss.Storyline_update_chapter, &ss.Timeline_update_chapter, &ss.Group_connection_update_chapter)
	if err != nil {
		return nil, err
	}

	return &ss, nil
}
