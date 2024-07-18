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
func PapperInfo(C *gin.Context) {
	rootPapperInfo, err := utils.GetRootPapperInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting Root papper info. " + err.Error()})
		return
	}
	papperInfo, err := utils.GetPapperInfo(C)

	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting papper info. " + err.Error()})
		return
	}

	query := "SELECT name, description, path, created_at, root_papper_id FROM pappers WHERE root_papper_id = $1 and id = $2"
	var papper models.Papper
	row := db.DB.QueryRow(query, rootPapperInfo.Id, papperInfo.ID)
	err = row.Scan(&papper.Name, &papper.Description, &papper.Path, &papper.Created_at, &papper.Root_papper_id)
	if err != nil {
		if err == sql.ErrNoRows {
			C.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "root folder not found in google drive. " + err.Error()})
			return
		}

		C.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	C.Set("papperInfo", papperInfo)

	C.Next()
}
