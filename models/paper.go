package models

import (
	"database/sql"
	"fmt"
	"time"

	"PaperTrail-fm.com/db"
	"google.golang.org/api/drive/v3"
)

type Paper struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Path        string    `json:"path"`
	Created_at  time.Time `json:"created_at"`
	World_id    string    `json:"world_id"`
	Color       string    `json:"color"`
	Order       int       `json:"order,omitempty"` // Use ponteiro para lidar com valores nulos
}

func (e *Paper) Save() error {
	var PaperId int

	// Verifica se o paper já existe no banco de dados
	query := `SELECT id FROM Papers WHERE name = $1 AND world_id = $2`
	err := db.DB.QueryRow(query, e.Name, e.World_id).Scan(&PaperId)

	if err != nil && err != sql.ErrNoRows {
		// Se ocorrer um erro diferente de "sem linhas encontradas", retorna o erro
		return fmt.Errorf("error checking paper existence: %v", err)
	}

	if err == sql.ErrNoRows {
		var maxOrder *int
		orderQuery := `SELECT MAX("order") FROM Papers WHERE world_id = $1`
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
		INSERT INTO Papers(id, name, description, path, created_at, world_id, "order", color) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
		_, err = db.DB.Exec(insertQuery, e.ID, e.Name, e.Description, e.Path, e.Created_at, e.World_id, newOrder, e.Color)
		if err != nil {
			return fmt.Errorf("error inserting paper: %v", err)
		}

		fmt.Println("Inserted new paper into database")
	} else {
		fmt.Println("Paper already exists in database")
	}

	return nil
}

func GetAllPapers() ([]Paper, error) {
	query := "SELECT * FROM Papers"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var Papers []Paper

	for rows.Next() {
		var paper Paper
		err := rows.Scan(&paper.ID, &paper.Name, &paper.Description, &paper.Path, &paper.Created_at, &paper.World_id, &paper.Color)

		if err != nil {
			return nil, err
		}

		Papers = append(Papers, paper)
	}

	return Papers, nil
}

func GetPaperByID(id int64) (*Paper, error) {
	query := "SELECT * FROM Papers WHERE id = ?"
	row := db.DB.QueryRow(query, id)

	var paper Paper
	err := row.Scan(&paper.ID, &paper.Name, &paper.Description, &paper.Path, &paper.Created_at, &paper.World_id, &paper.Color)
	if err != nil {
		return nil, err
	}

	return &paper, nil
}

func (paper *Paper) Update() error {

	query := `
	UPDATE Papers
	SET name = $1, description = $2, "order" = $3, color = $4
	WHERE id = $5
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(paper.Name, paper.Description, paper.Order, paper.Color, paper.ID)
	return err
}

func (paper *Paper) Delete() error {
	query := "DELETE FROM Papers WHERE id = ?"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(paper.ID)
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

func (paper *Paper) GetChapterList(driver *drive.Service, userAccessToken string) ([]chapterWithRevisions, error) {
	query := `SELECT id, name, description, created_at, Paper_id, "order" FROM chapters WHERE Paper_id = $1`
	rows, err := db.DB.Query(query, paper.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []chapterWithRevisions
	for rows.Next() {
		var chapter chapterWithRevisions
		if err := rows.Scan(&chapter.Id, &chapter.Name, &chapter.Description, &chapter.CreatedAt, &chapter.PaperID, &chapter.Order); err != nil {
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

func (paper *Paper) Get() error {
	query := `SELECT id, name, description, created_at, world_id, "order", color FROM Papers WHERE id = $1`
	row := db.DB.QueryRow(query, paper.ID)

	err := row.Scan(&paper.ID, &paper.Name, &paper.Description, &paper.Created_at, &paper.World_id, &paper.Order, &paper.Color)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no paper found with id %s", paper.ID)
		}
		return fmt.Errorf("error retrieving paper: %v", err)
	}
	return nil
}

func GetPaperListByWorldId(worldId string) ([]Paper, error) {
	Papers := []Paper{}

	PaperQuery := `
	SELECT id, name, description, created_at, "order", color
	FROM Papers
	WHERE world_id = $1
	`
	rows, err := db.DB.Query(PaperQuery, worldId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var paper Paper
		if err := rows.Scan(&paper.ID, &paper.Name, &paper.Description, &paper.Created_at, &paper.Order, &paper.Color); err != nil {
			return nil, err
		}

		Papers = append(Papers, paper)
	}
	return Papers, err
}
