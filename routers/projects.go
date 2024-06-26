package routes

import (
	"fmt"
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
	// ghClient.DeleteRepo(repoName)

	repo, err := ghClient.CreateRepo(repoName, papper.Description, true)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar repositório. " + err.Error()})
		return
	}
	fmt.Println("Repositório criado:", repo.GetHTMLURL())

	readmePath := "README.md"
	readmeContent := "# " + papper.Name + "\n\n" + papper.Description
	err = ghClient.CreateOrUpdateFile(repoName, userInfo.Email, userInfo.Name, readmePath, "initial README.md file", readmeContent)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar ou atualizar o arquivo. " + readmePath + err.Error()})
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

func GetFileUpdateList(context *gin.Context) {
	var file fileStr
	err := context.ShouldBindJSON(&file)

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
	repoName := utils.FormatRepositoryName(strconv.FormatInt(userInfo.ID, 10) + "_" + file.Papper)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	ghClient := githubclient.NewGitHubClient()
	commits, err := ghClient.ListCommitsOfFile(repoName, file.Path)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "file or project error." + err.Error()})
		return
	}
	context.JSON(http.StatusOK, commits)

}

func GetCommitDiff(context *gin.Context) {

	repo := context.Query("repo")
	sha := context.Query("sha")
	path := context.Query("path")
	ghClient := githubclient.NewGitHubClient()
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

	repoName := utils.FormatRepositoryName(strconv.FormatInt(userInfo.ID, 10) + "_" + repo)

	commits, err := ghClient.GetCommitDiff(repoName, sha, path)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, commits)

}
func GetFile(context *gin.Context) {
	repo := context.Query("repo")
	sha := context.Query("sha")
	path := context.Query("path")

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
	repoName := utils.FormatRepositoryName(strconv.FormatInt(userInfo.ID, 10) + "_" + repo)

	if repo == "" || sha == "" || path == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Missing required query parameters"})
		return
	}
	ghClient := githubclient.NewGitHubClient()

	content, err := ghClient.GetFileFromCommit(repoName, sha, path)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"content": content})

}

type fileStr struct {
	Papper string `json:"papper"`
	Path   string `json:"path"`
}

func Test() {

	ghClient := githubclient.NewGitHubClient()
	repoName := utils.FormatRepositoryName(strconv.FormatInt(1, 10) + "_" + "Um Belo livro para ficar rico")
	// ghClient.DeleteRepo(repoName)

	readmePath := "README.md"
	readmeContent := "# documento alterado para fins de teste"
	_ = ghClient.CreateOrUpdateFile(repoName, "email@email.com", "gabriel", readmePath, "commit que sera exibido . :)", readmeContent)

	fmt.Printf("Arquivo criado ou atualizado com sucesso (%s)\n", readmePath)

}
