package models

import (
	"database/sql"
	"fmt"
	"time"

	"PaperTrail-fm.com/db"
	"google.golang.org/api/drive/v3"
)

type Papper struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Path           string    `json:"path"`
	Created_at     time.Time `json:"created_at"`
	Root_papper_id string    `json:"root_papper_id"`
}

func (e *Papper) Save() error {
	var papperId int

	// Verifica se o papper já existe no banco de dados
	query := `SELECT id FROM pappers WHERE name = $1 AND root_papper_id = $2`
	err := db.DB.QueryRow(query, e.Name, e.Root_papper_id).Scan(&papperId)

	if err != nil && err != sql.ErrNoRows {
		// Se ocorrer um erro diferente de "sem linhas encontradas", retorna o erro
		return fmt.Errorf("error checking papper existence: %v", err)
	}

	if err == sql.ErrNoRows {
		// Se não há linhas (papper não existe), insere um novo registro
		insertQuery := `
		INSERT INTO pappers(id, name, description, path, created_at, root_papper_id) 
		VALUES ($1, $2, $3, $4, $5, $6)`
		_, err := db.DB.Exec(insertQuery, e.ID, e.Name, e.Description, e.Path, e.Created_at, e.Root_papper_id)
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
		err := rows.Scan(&papper.ID, &papper.Name, &papper.Description, &papper.Path, &papper.Created_at, &papper.Root_papper_id)

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
	err := row.Scan(&papper.ID, &papper.Name, &papper.Description, &papper.Path, &papper.Created_at, &papper.Root_papper_id)
	if err != nil {
		return nil, err
	}

	return &papper, nil
}

func (papper *Papper) Update() error {
	query := `
	UPDATE pappers
	SET name = ?, description = ?, path = ?, created_at = ?
	WHERE id = ?
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(papper.Name, papper.Description, papper.Path, papper.Created_at, papper.ID)
	return err
}

func (papper *Papper) Delete() error {
	query := "DELETE FROM pappers WHERE id = ?"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(papper.ID)
	return err
}

type chapterWithRevisions struct {
	Revisions []revision

	Chapter
	Iframe string
}

type revision struct {
	Id           string
	ModifiedTime string
}

func (papper *Papper) GetChapterList(driver *drive.Service, userAccessToken string) ([]chapterWithRevisions, error) {
	query := "SELECT id, name, description, created_at, Papper_id FROM chapters WHERE papper_id = $1"
	rows, err := db.DB.Query(query, papper.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []chapterWithRevisions
	for rows.Next() {
		var chapter chapterWithRevisions
		if err := rows.Scan(&chapter.Id, &chapter.Name, &chapter.Description, &chapter.CreatedAt, &chapter.PapperID); err != nil {
			return nil, err
		}
		revisions, _ := driver.Revisions.List(chapter.Id).Do()

		for _, element := range revisions.Revisions {

			var rev revision
			rev.Id = element.Id
			rev.ModifiedTime = element.ModifiedTime
			chapter.Revisions = append(chapter.Revisions, rev)
		}

		chapter.Iframe = GenerateSecureIframeURL(chapter.Id, userAccessToken)
		list = append(list, chapter)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func GenerateSecureIframeURL(fileID, token string) string {
	return fmt.Sprintf("https://docs.google.com/document/d/%s/edit?access_token=%s", fileID, token)
}
