package models

import (
	"database/sql"
	"fmt"
	"time"

	"PaperTrail-fm.com/db"
)

type RootPapper struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	Created_at time.Time `json:"created_at"`
	UserID     string    `json:"userId"`
}

func (rp *RootPapper) Save() error {
	var rpId string

	// Verifica se o papper já existe no banco de dados
	// query := `SELECT name, created_at FROM rootpappers WHERE id = $1`
	query := `SELECT name, created_at FROM rootpappers WHERE user_id = $1`
	err := db.DB.QueryRow(query, rp.UserID).Scan(&rpId)

	if err != nil && err != sql.ErrNoRows {
		// Se ocorrer um erro diferente de "sem linhas encontradas", retorna o erro
		return fmt.Errorf("error checking rootpappers existence: %v", err)
	}

	if err == sql.ErrNoRows {
		// Se não há linhas (papper não existe), insere um novo registro
		insertQuery := `
		INSERT INTO rootpappers(id, name, created_at, user_id) 
		VALUES ($1, $2, $3, $4)`
		_, err := db.DB.Exec(insertQuery, rp.Id, rp.Name, rp.Created_at, rp.UserID)
		if err != nil {
			return fmt.Errorf("error inserting papper: %v", err)
		}

		fmt.Println("Inserted new papper into database")
	} else {
		fmt.Println("Papper already exists in database")
	}

	return nil
}

func (rp *RootPapper) GetPapperList() ([]Papper, error) {
	query := "SELECT id, name, description, created_at FROM pappers WHERE root_papper_id = $1"
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
