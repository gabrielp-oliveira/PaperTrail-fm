package routes

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"PaperTrail-fm.com/db"
	"PaperTrail-fm.com/githubclient"
	"PaperTrail-fm.com/middlewares"
	"PaperTrail-fm.com/models"
	"PaperTrail-fm.com/utils"

	"github.com/gin-gonic/gin"
)

type basicPapperStructure struct {
	Description string `json:"description" binding:"required"`
	Name        string `json:"name" binding:"required"`
}

func GetAllPappers(context *gin.Context) {
	userId := context.GetInt64("userId")

	query := "SELECT * FROM pappers WHERE UserID = ?"
	rows, err := db.DB.Query(query, userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
		return
	}
	defer rows.Close()

	var pappers []models.Papper
	for rows.Next() {
		var papper models.Papper
		if err := rows.Scan(&papper.ID, &papper.Name, &papper.Description, &papper.UserID); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
			return
		}
		pappers = append(pappers, papper)
	}

	context.JSON(http.StatusOK, pappers)
}

func CreatePapper(context *gin.Context) {
	var basicPapper basicPapperStructure
	if err := context.ShouldBindJSON(&basicPapper); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	userInfoInterface, exists := context.Get("userInfo")
	if !exists {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível recuperar as informações do usuário"})
		return
	}

	userInfo, ok := userInfoInterface.(middlewares.UserBasicInfo)
	if !ok {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível converter as informações do usuário"})
		return
	}

	papper := models.Papper{
		Name:        basicPapper.Name,
		Description: basicPapper.Description,
		UserID:      userInfo.ID,
	}

	ghClient := githubclient.NewGitHubClient()
	repoName := utils.FormatRepositoryName(strconv.FormatInt(userInfo.ID, 10) + "_" + papper.Name)
	ghClient.DeleteRepo(repoName)

	repo, err := ghClient.CreateRepo(repoName, papper.Description, true)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar repositório"})
		log.Fatalf("Erro ao criar repositório: %v", err)
		return
	}
	fmt.Println("Repositório criado:", repo.GetHTMLURL())

	readmePath := "README.md"
	readmeContent := "# " + papper.Name + "\n\n" + papper.Description
	err = ghClient.CreateOrUpdateFile(repoName, userInfo.Email, userInfo.Name, readmePath, "initial README.md file", readmeContent)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar ou atualizar o arquivo"})
		log.Fatalf("Erro ao criar ou atualizar o arquivo (%s): %v", readmePath, err)
		return
	}
	fmt.Printf("Arquivo criado ou atualizado com sucesso (%s)\n", readmePath)
	err = papper.Save()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save paper"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Paper created successfully"})
}
