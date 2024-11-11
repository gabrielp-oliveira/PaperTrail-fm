package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"PaperTrail-fm.com/db"
	"github.com/google/uuid"
)

type StoryLine struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Created_at  time.Time `json:"created_at"`
	WorldID     string    `json:"world_id"`
	Order       int       `json:"order"`
}

func (t *StoryLine) Save() error {
	var storyLinesID string

	// Verifica se a storyLines já existe no banco de dados
	query := `SELECT id FROM storyLines WHERE id = $1`
	err := db.DB.QueryRow(query, t.Id).Scan(&storyLinesID)

	if err != nil && err != sql.ErrNoRows {
		// Se ocorrer um erro diferente de "sem linhas encontradas", retorna o erro
		return fmt.Errorf("error checking storyLines existence: %v", err)
	}

	if err == sql.ErrNoRows {
		var maxOrder *int
		orderQuery := `SELECT MAX("order") FROM storyLines WHERE world_id = $1 `
		err = db.DB.QueryRow(orderQuery, t.WorldID).Scan(&maxOrder)
		if err != nil {
			return fmt.Errorf("error getting max order: %v", err)
		}

		newOrder := 1
		if maxOrder != nil {
			newOrder = *maxOrder + 1
		} else {
			newOrder = 1
		}

		t.Order = newOrder

		// Se não há linhas (storyLines não existe), insere um novo registro
		insertQuery := `
		INSERT INTO storyLines(id, name, description, created_at, world_id, "order") 
		VALUES ($1, $2, $3, $4, $5, $6)`
		_, err := db.DB.Exec(insertQuery, t.Id, t.Name, t.Description, t.Created_at, t.WorldID, t.Order)
		if err != nil {
			return fmt.Errorf("error inserting storyLines: %v", err)
		}

		fmt.Println("Inserted new storyLines into database")
	} else {
		fmt.Println("storyLines already exists in database")
	}

	return nil
}

func GetstoryLinesByWorldId(worldId string) ([]StoryLine, error) {
	storyLines := []StoryLine{}

	storylinesQuery := `
		SELECT id, name, "order", description
		FROM storylines
		WHERE world_id = $1 ORDER BY "order"
	`
	rows, err := db.DB.Query(storylinesQuery, worldId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var storyLine StoryLine
		if err := rows.Scan(&storyLine.Id, &storyLine.Name, &storyLine.Order, &storyLine.Description); err != nil {
			return nil, err
		}
		storyLines = append(storyLines, storyLine)
	}

	return storyLines, err
}

func (t *StoryLine) Delete() error {
	// Verificar se o StoryLine existe
	query := `SELECT id FROM storyLines WHERE id = $1`
	err := db.DB.QueryRow(query, t.Id).Scan(&t.Id)
	if err == sql.ErrNoRows {
		return fmt.Errorf("StoryLine com id %s não existe", t.Id)
	} else if err != nil {
		return fmt.Errorf("Erro ao verificar a existência do StoryLine: %v", err)
	}

	// Encontrar IDs dos capítulos associados ao StoryLine
	query = `SELECT id FROM chapters WHERE storyline_id = $1`
	rows, err := db.DB.Query(query, t.Id)
	if err != nil {
		return fmt.Errorf("Erro ao buscar capítulos associados: %v", err)
	}
	defer rows.Close()

	var chapterIDs []string
	for rows.Next() {
		var chapterID string
		if err := rows.Scan(&chapterID); err != nil {
			return fmt.Errorf("Erro ao escanear capítulo: %v", err)
		}
		chapterIDs = append(chapterIDs, chapterID)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("Erro ao iterar sobre capítulos: %v", err)
	}

	// Se há capítulos, proceder com a atualização para remover os IDs de StoryLine, Timeline e Event
	if len(chapterIDs) > 0 {
		// Convertendo lista de IDs em uma string formatada para uso no SQL IN
		chapterIDPlaceholders := "'" + strings.Join(chapterIDs, "', '") + "'"

		// Atualizar capítulos para remover o storyline_id, timeline_id e event_id
		query = fmt.Sprintf(`UPDATE chapters SET storyline_id = NULL, timeline_id = NULL, event_id = NULL WHERE id IN (%s)`, chapterIDPlaceholders)
		if _, err := db.DB.Exec(query); err != nil {
			return fmt.Errorf("Erro ao remover IDs de StoryLine, Timeline e Event dos capítulos: %v", err)
		}
	}

	// Remover o próprio StoryLine
	query = `DELETE FROM storyLines WHERE id = $1`
	if _, err := db.DB.Exec(query, t.Id); err != nil {
		return fmt.Errorf("Erro ao remover StoryLine: %v", err)
	}

	fmt.Println("StoryLine removido e IDs de capítulos atualizados com sucesso.")
	return nil
}

func (t StoryLine) CreateBasicStoryLines(wiD string) ([]StoryLine, error) {
	var storyLinesList []StoryLine

	var str1 StoryLine
	str1.Name = "Main"
	str1.Description = "Main story line"
	str1.WorldID = wiD
	str1.Order = 1
	id := uuid.New().String()
	str1.Id = id
	str1.Created_at = time.Now()
	storyLinesList = append(storyLinesList, str1)
	err := str1.Save()
	if err != nil {
		return nil, fmt.Errorf("error removing storyLines: %v", err)
	}

	var str2 StoryLine
	str2.Name = "Seccondary"
	str2.Description = "Seccondary story line"
	str2.WorldID = wiD
	str2.Order = 2
	id = uuid.New().String()
	str2.Created_at = time.Now()
	str2.Id = id
	err = str2.Save()
	if err != nil {
		return nil, fmt.Errorf("error removing storyLines: %v", err)
	}
	storyLinesList = append(storyLinesList, str2)
	var str3 StoryLine
	str3.Name = "extra"
	str3.Description = "extra story line"
	str3.WorldID = wiD
	str3.Order = 3
	id = uuid.New().String()
	str3.Id = id
	str3.Created_at = time.Now()
	err = str3.Save()
	if err != nil {
		return nil, fmt.Errorf("error removing storyLines: %v", err)
	}
	storyLinesList = append(storyLinesList, str3)
	return storyLinesList, nil

}

func (t *StoryLine) Update() error {
	// id, name, description, range, "order"

	query := `
	UPDATE storylines
	SET name = $1, description = $2, "order" = $3
	WHERE id = $4
	`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(t.Name, t.Description, t.Order, t.Id)
	return err
}
