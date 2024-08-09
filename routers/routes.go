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

	authenticatedWorlds := authenticated.Group("/").Use(middlewares.WorldInfo)
	authenticatedWorlds.POST("/createPapper", CreatePapper)
	authenticatedWorlds.GET("/getPapperList", getPapperList)

	authenticatedWorlds.GET("/getRootData", GetRootData)

	authEvent := authenticated.Group("/").Use(middlewares.EventHandler)
	authEvent.POST("/insertEvent", InsertEvent)
	authEvent.DELETE("/removeEvent", RemoveEvent)
	authEvent.PUT("/updateEvent", UpdateEvent)

	authConnection := authenticated.Group("/").Use(middlewares.ConnectionHandler)
	authConnection.POST("/insertConnection", InsertConnection)
	authConnection.POST("/removeConnection", RemoveConnection)

	authenticatedPapper := authenticated.Group("/").Use(middlewares.PapperInfo)
	authenticatedPapper.POST("/createChapter", CreateChapter)
	authenticatedPapper.GET("/getChapterList", GetChapterList)
	authenticatedPapper.GET("/getChapter", GetChapter)

}
