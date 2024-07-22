package routes

import (
	"PaperTrail-fm.com/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)

	authenticated.POST("/createRootPapper", CreateRootPapper)
	authenticated.GET("/getRootPapperList", GetRootPapperList)

	authenticatedRootPapper := authenticated.Group("/").Use(middlewares.RootPapperInfo)
	authenticatedRootPapper.POST("/createPapper", CreatePapper)
	authenticatedRootPapper.GET("/getPapperList", getPapperList)

	authenticatedRootPapper.GET("/getRootData", GetRootData)

	authEvent := authenticated.Group("/").Use(middlewares.EventHandler)
	authEvent.POST("/insertEvent", InsertEvent)
	authEvent.DELETE("/removeEvent", RemoveEvent)

	authConnection := authenticated.Group("/").Use(middlewares.ConnectionHandler)
	authConnection.POST("/insertConnection", InsertConnection)
	authConnection.POST("/removeConnection", RemoveConnection)

	authenticatedPapper := authenticated.Group("/").Use(middlewares.PapperInfo)
	authenticatedPapper.POST("/createChapter", CreateChapter)
	authenticatedPapper.GET("/getChapterList", GetChapterList)
	authenticatedPapper.GET("/getChapter", GetChapter)

}
