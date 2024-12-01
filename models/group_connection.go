package models

import (
	"database/sql"
	"errors"
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

	query := `SELECT id FROM group_connections WHERE name = $1`
	err := db.DB.QueryRow(query, gc.Name).Scan(&id)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking connection existence: %v", err)
	}

	if err == sql.ErrNoRows {
		query := `
		INSERT INTO group_connections (id, name, description, color, world_id)
		VALUES ($1, $2, $3, $4, $5)`
		_, err := db.DB.Exec(query, gc.Id, gc.Name, gc.Description, gc.Color, gc.World_id)
		if err != nil {
			return fmt.Errorf("error inserting group_connections: %v", err)
		}

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

func (gc *GroupConnection) Delete() error {

	query := `SELECT world_id FROM group_connections WHERE id = $1`
	err := db.DB.QueryRow(query, gc.Id).Scan(&gc.World_id)

	if err != nil && err != sql.ErrNoRows {
		return errors.New("error checking connection existence: " + err.Error())
	}
	query = `DELETE FROM group_connections WHERE id = $1;`
	_, err = db.DB.Exec(query, gc.Id)
	if err != nil {
		return errors.New("error removing connection: " + err.Error())
	}

	query = `UPDATE connections SET group_id = $1 WHERE group_id = $2`
	var id = ""
	_, err = db.DB.Exec(query, id, gc.Id)
	if err != nil {
		return fmt.Errorf("error updating connections: %v", err)
	}

	return nil
}
func (gc *GroupConnection) Update() error {

	query := `
	update group_connections 
	SET name = $1, description = $2, color = $3 WHERE id = $4`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(gc.Name, gc.Description, gc.Color, gc.Id)
	return err
}
