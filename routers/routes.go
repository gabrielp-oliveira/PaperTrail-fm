package routes

import (
	"PaperTrail-fm.com/middlewares"
	"PaperTrail-fm.com/models"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)
	// authenticated.GET("/GetAllPappers", GetAllPappers)

	authenticated.POST("/createPapper", CreatePapper)
	authenticated.POST("/GetFileUpdateList", GetFileUpdateList)
	authenticated.GET("/GetCommitDiff", GetCommitDiff)
	authenticated.GET("/getFile", GetFile)

	server.GET("/Upload", models.Upload)
	server.GET("/download", models.Download)
}
