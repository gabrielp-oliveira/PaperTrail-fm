package middlewares

import (
	"database/sql"
	"net/http"
	"time"

	"PaperTrail-fm.com/db"
	"PaperTrail-fm.com/models"
	"PaperTrail-fm.com/utils"

	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")

	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
		return
	}

	userEmail, err := utils.VerifyToken(token)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
		return
	}
	query := "SELECT name, email, id, name, accessToken, refresh_token, token_expiry FROM users WHERE email = $1"
	row := db.DB.QueryRow(query, userEmail)

	var userInfo models.User
	err = row.Scan(&userInfo.Name, &userInfo.Email, &userInfo.ID, &userInfo.Name, &userInfo.AccessToken, &userInfo.RefreshToken, &userInfo.TokenExpiry)

	if err != nil {
		if err == sql.ErrNoRows {
			// Tratar caso não haja registros correspondentes
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
			return
		}
		// Outros erros podem ser tratados aqui

		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
		return
	}

	context.Set("userInfo", userInfo)

	context.Next()
}

type UserBasicInfo struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email" binding:"required"`
	Name         string    `json:"name"`
	Password     string    `json:"password"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refresh_token"`
	TokenExpiry  time.Time `json:"token_expiry"`
}
