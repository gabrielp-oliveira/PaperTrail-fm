package middlewares

import (
	"net/http"

	"PaperTrail-fm.com/db"
	"PaperTrail-fm.com/utils"

	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")

	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
		return
	}

	userId, err := utils.VerifyToken(token)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
		return
	}

	query := "SELECT name, email, id FROM users WHERE id = $1"
	row := db.DB.QueryRow(query, userId)

	var userInfo UserBasicInfo
	err = row.Scan(&userInfo.Name, &userInfo.Email, &userInfo.ID)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
	}

	context.Set("userInfo", userInfo)

	context.Next()
}

type UserBasicInfo struct {
	Email string `json:"email" binding:"required"`
	Name  string `json:"name" binding:"required"`
	ID    int64  `json:"id" binding:"required"`
}
