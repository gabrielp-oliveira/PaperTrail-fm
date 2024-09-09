package routes

import (
	"PaperTrail-fm.com/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)

	authenticated.POST("/createWorld", CreateWorld)
	authenticated.GET("/getWorldsList", GetWorldsList)
	authenticated.GET("/world", GetWorldData)

	authenticatedWorlds := authenticated.Group("/").Use(middlewares.WorldInfo)
	authenticatedWorlds.POST("/createPapper", CreatePapper)
	authenticatedWorlds.PUT("/updatePapper", UpdatePapper)
	authenticatedWorlds.GET("/getPapperList", getPapperList)

	// authenticatedWorlds.GET("/world", GetWorldData)

	authEvent := authenticated.Group("/").Use(middlewares.EventHandler)
	authEvent.POST("/insertEvent", InsertEvent)
	authEvent.DELETE("/removeEvent", RemoveEvent)
	authEvent.PUT("/updateEvent", UpdateEvent)

	authConnection := authenticated.Group("/").Use(middlewares.ConnectionHandler)
	authConnection.POST("/createConnection", CreateConnection)
	authConnection.POST("/removeConnection", RemoveConnection)
	// authenticatedPapper.GET("/papper", GetPapper)

	authenticatedPapper := authenticated.Group("/").Use(middlewares.PapperInfo)
	authenticatedPapper.POST("/createChapter", CreateChapter)
	authenticatedPapper.PUT("/updateChapter", UpdateChapter)
	authenticatedPapper.GET("/getChapterList", GetChapterList)

	authTimeline := authenticated.Group("/").Use(middlewares.TimelineInfo)
	authTimeline.POST("/createTimeline", CreateTimeline)
	authTimeline.PUT("/updateTimeline", UpdateTimeline)
	authTimeline.DELETE("/deleteTimeline", DeleteTimeline)

	authStoryline := authenticated.Group("/").Use(middlewares.StorylineInfo)
	authStoryline.POST("/createStoryline", CreateStoryline)
	authStoryline.PUT("/updateStoryline", UpdateStoryline)
	authStoryline.DELETE("/deleteStoryline", DeleteStoryline)

	authenticated.GET("/chapterUrl", GetChapterUrl)
	authenticated.GET("/chapter", GetChapter)
	authenticated.GET("/papper", GetPapper)
	// use the right middleware later
}
