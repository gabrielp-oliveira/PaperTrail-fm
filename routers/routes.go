package routes

import (
	"PaperTrail-fm.com/middlewares"
	"PaperTrail-fm.com/models"

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

	authenticatedPapper := authenticated.Group("/").Use(middlewares.PapperInfo)
	authenticatedPapper.POST("/createChapter", CreateChapter)

	// authenticatedRPapper := authenticated.Use(middlewares.PapperInfo)

	// authenticated.POST("/GetFileUpdateList", GetFileUpdateList)
	// authenticated.GET("/GetCommitDiff", GetCommitDiff)
	// authenticated.GET("/getFile", GetFile)
	// authenticated.POST("/createFile", CreateFile)
	// authenticated.POST("/updateFile", UpdateFile)
	// authenticated.POST("/GetFileList", GetFileList)
	// authenticated.POST("/GetDocxIdByProject", GetDocxIdByProject)

	server.GET("/Upload", models.Upload)
	server.GET("/download", models.Download)
}
