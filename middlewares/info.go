package middlewares

import (
	"database/sql"
	"net/http"

	"PaperTrail-fm.com/db"
	"PaperTrail-fm.com/models"
	"PaperTrail-fm.com/utils"

	"github.com/gin-gonic/gin"
)

func RootPapperInfo(C *gin.Context) {
	userInfo, err := utils.GetUserInfo(C)
	if err == sql.ErrNoRows {
		C.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	var papper models.Papper
	C.ShouldBindJSON(&papper)
	rootPapper, err := RootPapper(userInfo.ID, papper.Root_papper_id)

	if err == sql.ErrNoRows {
		C.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "root folder not found in google drive. " + err.Error()})
		return
	}
	C.Set("rootPapper", rootPapper)
	C.Set("papper", papper)

	C.Next()
}
func RootPapper(userId, rootPapperId string) (models.RootPapper, error) {

	query := "SELECT name, id FROM rootpappers WHERE user_id = $1 and id = $2"
	var rootPapper models.RootPapper
	row := db.DB.QueryRow(query, userId, rootPapperId)
	err := row.Scan(&rootPapper.Name, &rootPapper.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			return rootPapper, err
		}
		return rootPapper, err
	}

	return rootPapper, nil
}

func PapperInfo(C *gin.Context) {
	userInfo, err := utils.GetUserInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting User info. " + err.Error()})
		return
	}
	var Chapter models.Chapter
	C.ShouldBindJSON(&Chapter)

	rootPapperInfo, err := RootPapper(userInfo.ID, Chapter.RootPapperID)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting Root papper info. " + err.Error()})
		return
	}

	query := "SELECT name, id, description, path, created_at, root_papper_id FROM pappers WHERE root_papper_id = $1 and id = $2"
	row := db.DB.QueryRow(query, rootPapperInfo.Id, Chapter.PapperID)
	var papper models.Papper
	err = row.Scan(&papper.Name, &papper.ID, &papper.Description, &papper.Path, &papper.Created_at, &papper.Root_papper_id)
	if err != nil {
		if err == sql.ErrNoRows {
			C.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "root folder not found in google drive. " + err.Error()})
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
	rootPapper, err := RootPapper(userInfo.ID, event.Root_papper_id)
	C.Set("rootPapper", rootPapper)

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
		RootPapperId string
	}
	var connection Cnn
	C.ShouldBindJSON(&connection)

	rootPapper, err := RootPapper(userInfo.ID, connection.RootPapperId)
	C.Set("rootPapper", rootPapper)
	C.Set("connection", connection)

	C.Next()
}
