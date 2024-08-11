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

func Authenticate(C *gin.Context) {
	token := C.Request.Header.Get("Authorization")

	if token == "" {
		C.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
		return
	}

	userEmail, err := utils.VerifyToken(token)
	if err != nil {
		C.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
		return
	}
	query := "SELECT name, email,base_folder, id, name, accessToken, refresh_token, token_expiry FROM users WHERE email = $1"
	row := db.DB.QueryRow(query, userEmail)

	var userInfo models.User
	err = row.Scan(&userInfo.Name, &userInfo.Email, &userInfo.Base_folder, &userInfo.ID, &userInfo.Name, &userInfo.AccessToken, &userInfo.RefreshToken, &userInfo.TokenExpiry)

	if err != nil {
		if err == sql.ErrNoRows {
			C.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
			return
		}

		C.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
		return
	}
	newToken, err := userInfo.UpdateOAuthToken()
	if err != nil {
		C.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "error generating access token.."})
		return
	}

	C.Writer.Header().Set("accessToken", newToken.AccessToken)
	C.Writer.Header().Set("expiry", newToken.Expiry.Format(time.RFC3339))
	C.Writer.Header().Set("Access-Control-Expose-Headers", "accessToken, expiry")

	C.Set("userInfo", userInfo)
	C.Next()
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
