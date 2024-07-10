package middlewares

import (
	"database/sql"
	"net/http"

	"PaperTrail-fm.com/db"
	"PaperTrail-fm.com/models"
	"PaperTrail-fm.com/utils"

	"github.com/gin-gonic/gin"
)

func PapperInfo(C *gin.Context) {
	userInfo, err := utils.GetUserInfo(C)

	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user info. " + err.Error()})
		return
	}

	query := "SELECT name, id FROM rootpappers WHERE user_id = $1"
	var rootPapper models.RootPapper
	row := db.DB.QueryRow(query, userInfo.ID)
	err = row.Scan(&rootPapper.Name, &rootPapper.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			C.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "root folder not found in google drive. " + err.Error()})
			return
		}

		C.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	C.Set("rootPapper", rootPapper)

	C.Next()
}
