package models

import (
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
	query := `
	INSERT INTO pappers(name, description, path, dateTime, user_id) 
	VALUES (?, ?, ?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(e.Name, e.Description, e.Path, e.DateTime, e.UserID)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	e.ID = id
	return err
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

func (e Papper) Register(userId int64) error {
	query := "INSERT INTO registrations(event_id, user_id) VALUES (?, ?)"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)

	return err
}

func (e Papper) CancelRegistration(userId int64) error {
	query := "DELETE FROM registrations WHERE event_id = ? AND user_id = ?"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userId)

	return err
}
