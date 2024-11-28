package models

import (
	"database/sql"
	"fmt"

	"PaperTrail-fm.com/db"
)

type DescriptionRules struct {
	Id        string `json:"id"`
	Resource  string `json:"resource"`
	Max_pages string `json:"max_pages"`
}
type Description struct {
	Id               string `json:"id"`
	Resource_id      string `json:"resource_id"`
	Resource_type    string `json:"resource_type"`
	Description_data string `json:"description_data"`
	DescriptionRules
}

func (d *Description) CreateInitialDescription(resource_type string) error {
	var desc Description
	desc.Resource_id = d.Id
	desc.Resource_type = resource_type

	err := desc.Save()
	if err != nil {
		return err
	}
	return nil
}
func (d *Description) Save() error {
	var id string

	query := `SELECT id FROM descriptions WHERE id = $1`
	err := db.DB.QueryRow(query, d.Id).Scan(&id)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking descriptions existence: %v", err)
	}

	if err == sql.ErrNoRows {
		query := `
		INSERT INTO descriptions (id, resource_id, resource_type,description_data)
		VALUES ($1, $2, $3, $4)
	`
		_, err := db.DB.Exec(query, d.Id, d.Resource_id, d.Resource_type, d.Description_data)
		if err != nil {
			return fmt.Errorf("error inserting descriptions: %v", err)
		}

		fmt.Println("Inserted new descriptions into database")
	} else {
		return fmt.Errorf("descriptions already exists in database")
	}

	return nil
}

func (d *Description) Get() error {
	query := `SELECT resource_id, resource_type,description_data FROM descriptions WHERE id = $1`
	err := db.DB.QueryRow(query, d.Id).Scan(&d.Resource_id, &d.Resource_type, &d.Description_data)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no description found with id %s", d.Id)
		}
		return fmt.Errorf("error retrieving description: %v", err)
	}
	query = `SELECT max_pages FROM description_rules WHERE resource = $1`
	err = db.DB.QueryRow(query, d.Resource_type).Scan(&d.Max_pages)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no description_rules found with id %s", d.Id)
		}
		return fmt.Errorf("error retrieving description_rules: %v", err)
	}
	return nil
}
func (d *Description) GetByResourceId() error {
	query := `SELECT resource_type,description_data FROM descriptions WHERE resource_id = $1`
	err := db.DB.QueryRow(query, d.Resource_id).Scan(&d.Resource_type, &d.Description_data)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no description found with id %s", d.Id)
		}
		return fmt.Errorf("error retrieving description: %v", err)
	}
	query = `SELECT max_pages FROM description_rules WHERE resource = $1`
	err = db.DB.QueryRow(query, d.Resource_type).Scan(&d.Max_pages)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no description_rules found with id %s", d.Id)
		}
		return fmt.Errorf("error retrieving description_rules: %v", err)
	}
	return nil
}
func (d *Description) Update() error {
	query := `
	UPDATE chapters
	SET, description_data = $1
	WHERE id = $2
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(d.Description_data, d.Id)
	return err
}

// func GetConnectionsGroupListByWorldId(worldId string) ([]GroupConnection, error) {

// 	groupConnections := []GroupConnection{}
// 	connectionsQuery := `
// 		SELECT id, name, description,color, world_id
// 		FROM group_connections
// 		WHERE world_id = $1
// 	`
// 	rows, err := db.DB.Query(connectionsQuery, worldId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var gcs GroupConnection
// 		if err := rows.Scan(&gcs.Id, &gcs.Name, &gcs.Description, &gcs.Color, &gcs.World_id); err != nil {
// 			return nil, err
// 		}
// 		groupConnections = append(groupConnections, gcs)
// 	}

// 	return groupConnections, err
// }

// func (gc *GroupConnection) Update() error {

// 	query := `
// 	update group_connections
// 	SET name = $1, description = $2, color = $3 WHERE id = $4`
// 	stmt, err := db.DB.Prepare(query)

// 	if err != nil {
// 		return err
// 	}
// 	defer stmt.Close()
// 	_, err = stmt.Exec(gc.Name, gc.Description, gc.Color, gc.Id)
// 	return err
// }
