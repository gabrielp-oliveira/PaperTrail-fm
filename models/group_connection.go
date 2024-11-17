package models

import (
	"database/sql"
	"fmt"

	"PaperTrail-fm.com/db"
)

type GroupConnection struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	World_id    string `json:"world_id"`
	Color       string `json:"color"`
}

func (gc *GroupConnection) Save() error {
	var id string

	// Verificar se a conexão já existe na tabela connections
	query := `SELECT id FROM group_connections WHERE name = $1`
	err := db.DB.QueryRow(query, gc.Id).Scan(&id)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking connection existence: %v", err)
	}

	if err == sql.ErrNoRows {
		// Inserir nova conexão na tabela connections
		query := `
		INSERT INTO group_connections (id, name, description,color, world_id)
		VALUES ($1, $2, $3, $4, $5)
	`
		_, err := db.DB.Exec(query, gc.Id, gc.Name, gc.Description, gc.Color, gc, gc.World_id)
		if err != nil {
			return fmt.Errorf("error inserting group_connections: %v", err)
		}

		fmt.Println("Inserted new group_connections into database")
	} else {
		return fmt.Errorf("group_connections already exists in database")
	}

	return nil
}

func GetConnectionsGroupListByWorldId(worldId string) ([]GroupConnection, error) {

	groupConnections := []GroupConnection{}
	connectionsQuery := `
		SELECT id, name, description,color, world_id
		FROM group_connections
		WHERE world_id = $1
	`
	rows, err := db.DB.Query(connectionsQuery, worldId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var gcs GroupConnection
		if err := rows.Scan(&gcs.Id, &gcs.Name, &gcs.Description, &gcs.Color, &gcs.World_id); err != nil {
			return nil, err
		}
		groupConnections = append(groupConnections, gcs)
	}

	return groupConnections, err
}
