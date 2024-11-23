package utils

import (
	"errors"

	"PaperTrail-fm.com/models"
	"github.com/gin-gonic/gin"
)

func GetWorldsInfo(context *gin.Context) (models.World, error) {
	worldsI, exists := context.Get("world")
	queryWorldsI := context.Query("id")

	var worldId string

	if !exists {
		if queryWorldsI == "" {
			return models.World{}, errors.New("unable to retrieve worlds information")
		} else {
			worldId = queryWorldsI
		}
	} else {
		// Considerando que worldsI é um tipo que possui um campo Id
		if world, ok := worldsI.(models.World); ok {
			worldId = world.Id
		} else {
			return models.World{}, errors.New("invalid world information")
		}
	}

	var worldInfo models.World
	worldInfo.Id = worldId
	err := worldInfo.Get()
	if err != nil {
		return models.World{}, errors.New("unable to retrieve worlds information")
	}

	return worldInfo, nil
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

func GetContextInfo[T any](context *gin.Context, key string) (T, error) {
	var result T

	value, exists := context.Get(key)
	if !exists {
		return result, errors.New("unable to retrieve information from context: " + key)
	}

	// Realiza o type assertion para o tipo `T`
	info, ok := value.(T)
	if !ok {
		return result, errors.New("unable to assert type for key: " + key)
	}

	return info, nil
}
