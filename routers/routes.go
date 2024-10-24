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
	authenticatedWorlds.POST("/createPaper", CreatePaper)
	authenticatedWorlds.PUT("/updatePaper", UpdatePaper)
	authenticated.PUT("/updatePaperList", UpdatePaperList)
	authenticatedWorlds.GET("/getPaperList", getPaperList)

	// authenticatedWorlds.GET("/world", GetWorldData)

	authEvent := authenticated.Group("/").Use(middlewares.EventHandler)
	authEvent.POST("/insertEvent", InsertEvent)
	authEvent.DELETE("/removeEvent", RemoveEvent)
	authEvent.PUT("/updateEvent", UpdateEvent)

	authConnection := authenticated.Group("/").Use(middlewares.ConnectionHandler)
	authConnection.POST("/createConnection", CreateConnection)
	authConnection.POST("/removeConnection", RemoveConnection)
	// authenticatedPaper.GET("/paper", GetPaper)

	authenticatedPaper := authenticated.Group("/").Use(middlewares.PaperInfo)
	authenticatedPaper.POST("/createChapter", CreateChapter)
	authenticatedPaper.PUT("/updateChapter", UpdateChapter)
	authenticated.PUT("/updateChapterList", UpdateChapterList)
	authenticatedPaper.GET("/getChapterList", GetChapterList)

	authTimeline := authenticated.Group("/").Use(middlewares.TimelineInfo)
	authTimeline.POST("/createTimeline", CreateTimeline)
	authTimeline.PUT("/updateTimeline", UpdateTimeline)
	authTimeline.DELETE("/deleteTimeline", DeleteTimeline)

	authStoryline := authenticated.Group("/").Use(middlewares.StorylineInfo)
	authStoryline.POST("/createStoryline", CreateStoryline)
	authStoryline.PUT("/updateStoryline", UpdateStoryline)
	authenticated.PUT("/updateStorylineList", UpdateStorylineList)
	authStoryline.DELETE("/deleteStoryline", DeleteStoryline)

	authenticated.GET("/chapterUrl", GetChapterUrl)
	authenticated.GET("/chapter", GetChapter)
	authenticated.GET("/paper", GetPaper)
	// use the right middleware later
}
