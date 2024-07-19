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
func GetChapterInfo(context *gin.Context) (models.Chapter, error) {

	ChapterI, exists := context.Get("chapter")
	if !exists {
		return models.Chapter{}, errors.New("Unable to retrieve user information")
	}

	ChapterInfo, ok := ChapterI.(models.Chapter)
	if !ok {
		return models.Chapter{}, errors.New("Unable to retrieve user information")
	}

	return ChapterInfo, nil
}
