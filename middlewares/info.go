package middlewares

import (
	"database/sql"
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
	var papper models.Papper
	C.ShouldBindJSON(&papper)
	world, err := World(userInfo.ID, papper.World_id)

	if err == sql.ErrNoRows {
		C.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "world  not found in google drive. " + err.Error()})
		return
	}
	C.Set("world", world)
	C.Set("papper", papper)

	C.Next()
}
func World(userId, worldsId string) (models.Worlds, error) {

	query := "SELECT name, id FROM worlds WHERE user_id = $1 and id = $2"
	var worlds models.Worlds
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

func PapperInfo(C *gin.Context) {
	userInfo, err := utils.GetUserInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting User info. " + err.Error()})
		return
	}
	var Chapter models.Chapter
	C.ShouldBindJSON(&Chapter)

	worldInfo, err := World(userInfo.ID, Chapter.WorldsID)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting Root papper info. " + err.Error()})
		return
	}

	query := "SELECT name, id, description, path, created_at, world_id FROM pappers WHERE world_id = $1 and id = $2"
	row := db.DB.QueryRow(query, worldInfo.Id, Chapter.PapperID)
	var papper models.Papper
	err = row.Scan(&papper.Name, &papper.ID, &papper.Description, &papper.Path, &papper.Created_at, &papper.World_id)
	if err != nil {
		if err == sql.ErrNoRows {
			C.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "world  not found in google drive. " + err.Error()})
			return
		}

		C.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	C.Set("papper", papper)
	C.Set("chapter", Chapter)

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
	type Cnn struct {
		models.Connection
		WorldsId string
	}
	var connection Cnn
	C.ShouldBindJSON(&connection)

	world, err := World(userInfo.ID, connection.WorldsId)
	C.Set("world", world)
	C.Set("connection", connection)

	C.Next()
}
