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
func GetPapperInfo(context *gin.Context) (models.Papper, error) {

	papperData, exists := context.Get("papper")
	if !exists {
		return models.Papper{}, errors.New("Unable to retrieve papper information")

	}

	// Fazer o cast para o tipo correto
	papper, ok := papperData.(models.Papper)
	if !ok {
		return models.Papper{}, errors.New("Unable to retrieve Papper information")

	}
	return papper, nil

}
