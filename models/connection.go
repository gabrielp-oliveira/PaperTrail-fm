package models

import (
	"database/sql"
	"errors"
	"fmt"

	"PaperTrail-fm.com/db"
)

type Connection struct {
	Id              string `json:"id"`
	SourceChapterID string `json:"sourceChapterID"`
	TargetChapterID string `json:"targetChapterID"`
	World_id        string `json:"world_id"`
}

func (cnn *Connection) Save() error {
	var id string

	// Verificar se a conexão já existe na tabela connections
	query := `SELECT id FROM connections WHERE source_chapter_id = $1 and  target_chapter_id = $2`
	err := db.DB.QueryRow(query, cnn.SourceChapterID, cnn.TargetChapterID).Scan(&id)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking connection existence: %v", err)
	}
	query = `SELECT id FROM connections WHERE  target_chapter_id = $1 and source_chapter_id = $2`
	err = db.DB.QueryRow(query, cnn.SourceChapterID, cnn.TargetChapterID).Scan(&id)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking connection existence: %v", err)
	}

	if err == sql.ErrNoRows {
		// Inserir nova conexão na tabela connections
		query := `
		INSERT INTO connections (id, source_chapter_id, target_chapter_id, world_id)
		VALUES ($1, $2, $3, $4)
	`
		_, err := db.DB.Exec(query, cnn.Id, cnn.SourceChapterID, cnn.TargetChapterID, cnn.World_id)
		if err != nil {
			return fmt.Errorf("error inserting connection: %v", err)
		}

		fmt.Println("Inserted new connection into database")
	} else {
		return fmt.Errorf("Connection already exists in database")
	}

	return nil
}

func (cnn *Connection) Delete() error {
	query := `SELECT id FROM connections WHERE id = $1`
	err := db.DB.QueryRow(query, cnn.Id).Scan(&cnn.Id)

	if err != nil && err != sql.ErrNoRows {
		return errors.New("error checking connection existence: " + err.Error())
	}
	query = `DELETE FROM connections WHERE id = $1;`
	_, err = db.DB.Exec(query, cnn.Id)
	if err != nil {
		return errors.New("error removing connection: " + err.Error())
	}

	fmt.Println("Removed connection from database")

	return nil
}

func GetConnectionsListByWorldId(worldId string) ([]Connection, error) {

	connections := []Connection{}
	connectionsQuery := `
		SELECT id, source_chapter_id, target_chapter_id
		FROM connections
		WHERE source_chapter_id IN (SELECT id FROM chapters WHERE world_id = $1)
	`
	rows, err := db.DB.Query(connectionsQuery, worldId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var connection Connection
		if err := rows.Scan(&connection.Id, &connection.SourceChapterID, &connection.TargetChapterID); err != nil {
			return nil, err
		}
		connection.World_id = worldId
		connections = append(connections, connection)
	}

	return connections, err
}
