package routes

import (
	"net/http"

	"PaperTrail-fm.com/db"

	"github.com/gin-gonic/gin"
)

func GetAllPappers(context *gin.Context) {
	userId := context.GetInt64("userId")

	query := "SELECT * FROM pappers WHERE UserID = ?"
	rows, err := db.DB.Query(query, userId)

	if err != nil {
		return
	}
	context.JSON(http.StatusOK, rows)

}
