package utils

import (
	"errors"

	"PaperTrail-fm.com/models"
	"github.com/gin-gonic/gin"
)

func GetWorldsInfo(context *gin.Context) (models.Worlds, error) {

	worldsI, exists := context.Get("world")
	if !exists {
		return models.Worlds{}, errors.New("unable to retrieve worlds information")
	}

	worldsInfo, ok := worldsI.(models.Worlds)
	if !ok {
		return models.Worlds{}, errors.New("unable to retrieve worlds information")
	}

	return worldsInfo, nil
}
func GetConnectionInfo(context *gin.Context) (models.Connection, error) {

	connection, exists := context.Get("connection")
	if !exists {
		return models.Connection{}, errors.New("unable to retrieve worlds information")
	}

	connectionnfo, ok := connection.(models.Connection)
	if !ok {
		return models.Connection{}, errors.New("unable to retrieve worlds information")
	}

	return connectionnfo, nil
}
func GetEventInfo(context *gin.Context) (models.Event, error) {

	event, exists := context.Get("event")
	if !exists {
		return models.Event{}, errors.New("unable to retrieve worlds information")
	}

	eventInfo, ok := event.(models.Event)
	if !ok {
		return models.Event{}, errors.New("unable to retrieve worlds information")
	}

	return eventInfo, nil
}

func GetChapterInfo(context *gin.Context) (models.Chapter, error) {

	ChapterI, exists := context.Get("chapter")
	if !exists {
		return models.Chapter{}, errors.New("unable to retrieve worlds information")
	}

	ChapterInfo, ok := ChapterI.(models.Chapter)
	if !ok {
		return models.Chapter{}, errors.New("unable to retrieve worlds information")
	}

	return ChapterInfo, nil
}
