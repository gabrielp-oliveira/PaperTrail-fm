package main

import (
	"log"

	routes "PaperTrail-fm.com/routers"

	"PaperTrail-fm.com/s3Instance"

	"PaperTrail-fm.com/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	s3Instance.InitS3()
	db.InitDB()
	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":9090") // localhost:8080
}
