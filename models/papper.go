package models

import (
	"database/sql"
	"fmt"
	"time"

	"PaperTrail-fm.com/db"
	"google.golang.org/api/drive/v3"
)

type Papper struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Path        string    `json:"path"`
	Created_at  time.Time `json:"created_at"`
	World_id    string    `json:"world_id"`
	Order       int       `json:"order,omitempty"` // Use ponteiro para lidar com valores nulos
}

func (e *Papper) Save() error {
	var papperId int

	// Verifica se o papper já existe no banco de dados
	query := `SELECT id FROM pappers WHERE name = $1 AND world_id = $2`
	err := db.DB.QueryRow(query, e.Name, e.World_id).Scan(&papperId)

	if err != nil && err != sql.ErrNoRows {
		// Se ocorrer um erro diferente de "sem linhas encontradas", retorna o erro
		return fmt.Errorf("error checking papper existence: %v", err)
	}

	if err == sql.ErrNoRows {
		var maxOrder *int
		orderQuery := `SELECT MAX("order") FROM pappers WHERE world_id = $1`
		err := db.DB.QueryRow(orderQuery, e.World_id).Scan(&maxOrder)
		if err != nil {
			return fmt.Errorf("error getting max order: %v", err)
		}

		newOrder := 1
		if maxOrder != nil {
			newOrder = *maxOrder + 1
		}

		e.Order = newOrder
		insertQuery := `
		INSERT INTO pappers(id, name, description, path, created_at, world_id, "order") 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
		_, err = db.DB.Exec(insertQuery, e.ID, e.Name, e.Description, e.Path, e.Created_at, e.World_id, newOrder)
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
		err := rows.Scan(&papper.ID, &papper.Name, &papper.Description, &papper.Path, &papper.Created_at, &papper.World_id)

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
	err := row.Scan(&papper.ID, &papper.Name, &papper.Description, &papper.Path, &papper.Created_at, &papper.World_id)
	if err != nil {
		return nil, err
	}

	return &papper, nil
}

func (papper *Papper) Update() error {
	query := `
	UPDATE pappers
	SET name = $1, description = $2, "order" = $3
	WHERE id = $4
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(papper.Name, papper.Description, papper.Order, papper.ID)
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
	query := `SELECT id, name, description, created_at, Papper_id, "order" FROM chapters WHERE papper_id = $1`
	rows, err := db.DB.Query(query, papper.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []chapterWithRevisions
	for rows.Next() {
		var chapter chapterWithRevisions
		if err := rows.Scan(&chapter.Id, &chapter.Name, &chapter.Description, &chapter.CreatedAt, &chapter.PapperID, &chapter.Order); err != nil {
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

func (papper *Papper) Get() error {
	query := `SELECT id, name, description, created_at, world_id, "order" FROM pappers WHERE id = $1`
	row := db.DB.QueryRow(query, papper.ID)

	err := row.Scan(&papper.ID, &papper.Name, &papper.Description, &papper.Created_at, &papper.World_id, &papper.Order)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no papper found with id %s", papper.ID)
		}
		return fmt.Errorf("error retrieving papper: %v", err)
	}
	return nil
}

func GetPapperListByWorldId(worldId string) ([]Papper, error) {
	pappers := []Papper{}

	papperQuery := `
	SELECT id, name, description, created_at, "order"
	FROM pappers
	WHERE world_id = $1
	`
	rows, err := db.DB.Query(papperQuery, worldId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var papper Papper
		if err := rows.Scan(&papper.ID, &papper.Name, &papper.Description, &papper.Created_at, &papper.Order); err != nil {
			return nil, err
		}

		pappers = append(pappers, papper)
	}
	return pappers, err
}
