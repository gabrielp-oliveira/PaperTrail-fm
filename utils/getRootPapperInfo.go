package utils

import (
	"errors"

	"PaperTrail-fm.com/models"
	"github.com/gin-gonic/gin"
)

func GetRootPapperInfo(context *gin.Context) (models.RootPapper, error) {

	rootPapperI, exists := context.Get("rootPapper")
	if !exists {
		return models.RootPapper{}, errors.New("Unable to retrieve user information")
	}

	rootPapperInfo, ok := rootPapperI.(models.RootPapper)
	if !ok {
		return models.RootPapper{}, errors.New("Unable to retrieve user information")
	}

	return rootPapperInfo, nil
}
