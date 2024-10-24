package utils

import (
	"errors"

	"PaperTrail-fm.com/models"
	"github.com/gin-gonic/gin"
)

func GetUserInfo(context *gin.Context) (models.User, error) {

	userInfoInterface, exists := context.Get("userInfo")
	if !exists {
		return models.User{}, errors.New("unable to retrieve user information")
	}

	userInfo, ok := userInfoInterface.(models.User)
	if !ok {
		return models.User{}, errors.New("unable to retrieve user information")
	}

	return userInfo, nil
}
func GetPaperInfo(context *gin.Context) (models.Paper, error) {

	PaperData, exists := context.Get("paper")
	if !exists {
		return models.Paper{}, errors.New("Unable to retrieve paper information")

	}

	// Fazer o cast para o tipo correto
	paper, ok := PaperData.(models.Paper)
	if !ok {
		return models.Paper{}, errors.New("Unable to retrieve Paper information")

	}
	return paper, nil

}
