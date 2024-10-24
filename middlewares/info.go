package middlewares

import (
	"database/sql"
	"fmt"
	"net/http"

	"PaperTrail-fm.com/db"
	"PaperTrail-fm.com/models"
	"PaperTrail-fm.com/utils"

	"github.com/gin-gonic/gin"
)

func WorldInfo(C *gin.Context) {
	userInfo, err := utils.GetUserInfo(C)
	if err == sql.ErrNoRows {
		C.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	world_id := C.Query("world_id")
	var paper models.Paper

	if world_id != "" {
		paper.World_id = world_id
	} else {
		C.ShouldBindJSON(&paper)
	}
	world, err := World(userInfo.ID, paper.World_id)

	if err == sql.ErrNoRows {
		C.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "world  not found in google drive. " + err.Error()})
		return
	}
	C.Set("world", world)
	C.Set("paper", paper)

	C.Next()
}
func World(userId, worldsId string) (models.World, error) {

	query := "SELECT name, id FROM worlds WHERE user_id = $1 and id = $2"
	var worlds models.World
	row := db.DB.QueryRow(query, userId, worldsId)
	err := row.Scan(&worlds.Name, &worlds.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return worlds, err
		}
		return worlds, err
	}

	return worlds, nil
}

func PaperInfo(C *gin.Context) {
	userInfo, err := utils.GetUserInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting User info. " + err.Error()})
		return
	}

	var Chapter models.Chapter
	C.ShouldBindJSON(&Chapter)

	worldInfo, err := World(userInfo.ID, Chapter.WorldsID)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting Root paper info. " + err.Error()})
		return
	}

	query := `SELECT name, id, description, path, created_at, world_id, "order" FROM Papers WHERE id = $1`
	row := db.DB.QueryRow(query, worldInfo.Id, Chapter.PaperID)
	var paper models.Paper
	err = row.Scan(&paper.Name, &paper.ID, &paper.Description, &paper.Path, &paper.Created_at, &paper.World_id, &paper.Order)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Errorf("paper not found.")
		}
	}

	C.Set("paper", paper)
	C.Set("chapter", Chapter)

	C.Next()
}
func StorylineInfo(C *gin.Context) {
	userInfo, err := utils.GetUserInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting User info. " + err.Error()})
		return
	}

	var stl models.StoryLine
	C.ShouldBindJSON(&stl)

	worldInfo, err := World(userInfo.ID, stl.WorldsID)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting Root paper info. " + err.Error()})
		return
	}

	query := "SELECT name, description, created_at, worldsId, 'order' FROM storyLines WHERE id = $1"
	row := db.DB.QueryRow(query, stl.Id)

	err = row.Scan(&stl.Name, &stl.Description, &stl.Created_at, &stl.WorldsID, &stl.Order)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Errorf("paper not found.")
		}
	}

	C.Set("world", worldInfo)
	C.Set("storyline", stl)

	C.Next()
}
func TimelineInfo(C *gin.Context) {
	userInfo, err := utils.GetUserInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting User info. " + err.Error()})
		return
	}

	var tl models.Timeline
	C.ShouldBindJSON(&tl)

	worldInfo, err := World(userInfo.ID, tl.WorldsID)
	if err != nil {
		fmt.Errorf("Error getting Root paper info. " + err.Error())
	}

	query := `SELECT id, name, description, range, worldsId, 'order' FROM Timelines WHERE id = $1 ORDER BY "order"`
	row := db.DB.QueryRow(query, tl.Id)

	err = row.Scan(&tl.Id, &tl.Name, &tl.Description, &tl.Range, &tl.WorldsID, &tl.Order)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Errorf("paper not found.")
		}
	}

	C.Set("world", worldInfo)
	C.Set("timeline", tl)

	C.Next()
}

func EventHandler(C *gin.Context) {
	userInfo, err := utils.GetUserInfo(C)
	if err == sql.ErrNoRows {
		C.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	var event models.Event
	C.ShouldBindJSON(&event)
	world, err := World(userInfo.ID, event.World_id)
	C.Set("world", world)

	C.Set("event", event)

	C.Next()
}
func ConnectionHandler(C *gin.Context) {
	userInfo, err := utils.GetUserInfo(C)
	if err == sql.ErrNoRows {
		C.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	var connection models.Connection
	C.ShouldBindJSON(&connection)

	world, err := World(userInfo.ID, connection.World_id)
	if err != nil {
		C.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	C.Set("world", world)
	C.Set("connection", connection)

	C.Next()
}
