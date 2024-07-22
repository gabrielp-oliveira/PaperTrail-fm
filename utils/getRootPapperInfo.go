package utils

import (
	"errors"

	"PaperTrail-fm.com/models"
	"github.com/gin-gonic/gin"
)

func GetRootPapperInfo(context *gin.Context) (models.RootPapper, error) {

	rootPapperI, exists := context.Get("rootPapper")
	if !exists {
		return models.RootPapper{}, errors.New("unable to retrieve rootPapper information")
	}

	rootPapperInfo, ok := rootPapperI.(models.RootPapper)
	if !ok {
		return models.RootPapper{}, errors.New("unable to retrieve rootPapper information")
	}

	return rootPapperInfo, nil
}
func GetConnectionInfo(context *gin.Context) (models.Connection, error) {

	connection, exists := context.Get("connection")
	if !exists {
		return models.Connection{}, errors.New("unable to retrieve rootPapper information")
	}

	connectionnfo, ok := connection.(models.Connection)
	if !ok {
		return models.Connection{}, errors.New("unable to retrieve rootPapper information")
	}

	return connectionnfo, nil
}
func GetEventInfo(context *gin.Context) (models.Event, error) {

	event, exists := context.Get("event")
	if !exists {
		return models.Event{}, errors.New("unable to retrieve rootPapper information")
	}

	eventInfo, ok := event.(models.Event)
	if !ok {
		return models.Event{}, errors.New("unable to retrieve rootPapper information")
	}

	return eventInfo, nil
}

func GetChapterInfo(context *gin.Context) (models.Chapter, error) {

	ChapterI, exists := context.Get("chapter")
	if !exists {
		return models.Chapter{}, errors.New("unable to retrieve rootPapper information")
	}

	ChapterInfo, ok := ChapterI.(models.Chapter)
	if !ok {
		return models.Chapter{}, errors.New("unable to retrieve rootPapper information")
	}

	return ChapterInfo, nil
}
