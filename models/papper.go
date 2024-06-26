package models

import (
	"database/sql"
	"fmt"
	"time"

	"PaperTrail-fm.com/db"
)

type Papper struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Path        string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int64     `binding:"required"`
}

func (e *Papper) Save() error {
	var papperId int

	// Verifica se o papper já existe no banco de dados
	query := `SELECT id FROM pappers WHERE name = $1 AND user_id = $2`
	err := db.DB.QueryRow(query, e.Name, e.UserID).Scan(&papperId)

	if err != nil && err != sql.ErrNoRows {
		// Se ocorrer um erro diferente de "sem linhas encontradas", retorna o erro
		return fmt.Errorf("error checking papper existence: %v", err)
	}

	if err == sql.ErrNoRows {
		// Se não há linhas (papper não existe), insere um novo registro
		insertQuery := `
		INSERT INTO pappers(name, description, path, dateTime, user_id) 
		VALUES ($1, $2, $3, $4, $5)`
		_, err := db.DB.Exec(insertQuery, e.Name, e.Description, e.Path, e.DateTime, e.UserID)
		if err != nil {
			return fmt.Errorf("error inserting papper: %v", err)
		}

		fmt.Println("Inserted new papper into database")
	} else {
		fmt.Println("Papper already exists in database")
	}

	return nil
}

func GetAllPappers() ([]Papper, error) {
	query := "SELECT * FROM pappers"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pappers []Papper

	for rows.Next() {
		var papper Papper
		err := rows.Scan(&papper.ID, &papper.Name, &papper.Description, &papper.Path, &papper.DateTime, &papper.UserID)

		if err != nil {
			return nil, err
		}

		pappers = append(pappers, papper)
	}

	return pappers, nil
}

func GetPapperByID(id int64) (*Papper, error) {
	query := "SELECT * FROM pappers WHERE id = ?"
	row := db.DB.QueryRow(query, id)

	var papper Papper
	err := row.Scan(&papper.ID, &papper.Name, &papper.Description, &papper.Path, &papper.DateTime, &papper.UserID)
	if err != nil {
		return nil, err
	}

	return &papper, nil
}

func (papper Papper) Update() error {
	query := `
	UPDATE pappers
	SET name = ?, description = ?, path = ?, dateTime = ?
	WHERE id = ?
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(papper.Name, papper.Description, papper.Path, papper.DateTime, papper.ID)
	return err
}

func (papper Papper) Delete() error {
	query := "DELETE FROM pappers WHERE id = ?"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(papper.ID)
	return err
}
