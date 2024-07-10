package utils

import (
	"errors"

	"PaperTrail-fm.com/models"
	"github.com/gin-gonic/gin"
)

func GetUserInfo(context *gin.Context) (models.User, error) {

	userInfoInterface, exists := context.Get("userInfo")
	if !exists {
		return models.User{}, errors.New("Unable to retrieve user information")
	}

	userInfo, ok := userInfoInterface.(models.User)
	if !ok {
		return models.User{}, errors.New("Unable to retrieve user information")
	}

	return userInfo, nil
}
