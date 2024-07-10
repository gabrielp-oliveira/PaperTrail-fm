package routes

import (
	"context"
	"fmt"
	"net/http"

	"PaperTrail-fm.com/githubclient"
	"PaperTrail-fm.com/googleClient"
	"PaperTrail-fm.com/models"
	"PaperTrail-fm.com/utils"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	"github.com/gin-gonic/gin"
)

type basicPapperStructure struct {
	Description string `json:"description" binding:"required"`
	Name        string `json:"name" binding:"required"`
}

var googleOauthConfig = googleClient.StartCredentials()

// func CreatePapper(context *gin.Context) {
// 	var basicPapper basicPapperStructure
// 	if err := context.ShouldBindJSON(&basicPapper); err != nil {
// 		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
// 		return
// 	}
// 	userInfo, err := utils.GetUserInfo(context)
// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	papper := models.Papper{
// 		Name:        basicPapper.Name,
// 		Description: basicPapper.Description,
// 		UserID:      userInfo.ID,
// 	}
// 	ghClient := githubclient.NewGitHubClient()
// 	repoName := utils.FormatRepositoryName(userInfo.ID + "_" + papper.Name)
// 	// ghClient.DeleteRepo(repoName)
// 	repo, err := ghClient.CreateRepo(repoName, papper.Description, true)
// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar repositório. " + err.Error()})
// 		return
// 	}
// 	fmt.Println("Repositório criado:", repo.GetHTMLURL())
// 	readmePath := "README.md"
// 	readmeContent := "# " + papper.Name + "\n\n" + papper.Description
// 	err = ghClient.CreateOrUpdateFile(repoName, userInfo.Email, userInfo.Name, readmePath, "initial README.md file", readmeContent)
// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar ou atualizar o arquivo. " + readmePath + err.Error()})
// 		return
// 	}
// 	fmt.Printf("Arquivo criado ou atualizado com sucesso (%s)\n", readmePath)
// 	err = papper.Save()
// 	if err != nil {
// 		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save paper"})
// 		return
// 	}
// 	context.JSON(http.StatusOK, gin.H{"message": "Paper created successfully"})
// }

func GetFileUpdateList(context *gin.Context) {
	var file fileStr
	err := context.ShouldBindJSON(&file)

	userInfo, err := utils.GetUserInfo(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	repoName := utils.FormatRepositoryName(userInfo.ID + "_" + file.Papper)

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

	userInfo, err := utils.GetUserInfo(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	repoName := utils.FormatRepositoryName(userInfo.ID + "_" + repo)

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

	userInfo, err := utils.GetUserInfo(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	repoName := utils.FormatRepositoryName(userInfo.ID + "_" + repo)

	if repo == "" || sha == "" || path == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Missing required query parameters"})
		return
	}
	ghClient := githubclient.NewGitHubClient()

	content, ext, err := ghClient.GetFileFromCommit(repoName, sha, path)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	context.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s%s", "file", ext))

	context.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", content)
}

func CreateFile(context *gin.Context) {

	var basicPapper fileStr
	if err := context.ShouldBindJSON(&basicPapper); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	ghClient := githubclient.NewGitHubClient()

	userInfo, err := utils.GetUserInfo(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	repoName := utils.FormatRepositoryName(userInfo.ID + "_" + basicPapper.Papper)
	Filename := utils.FormatRepositoryName(basicPapper.Path)

	err = ghClient.CreateFile(repoName, userInfo.Email, userInfo.Name, Filename+"/"+Filename, "initial File", "initial doc")

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "nao foi possivel criar arquivo. " + repoName + Filename})
		return
	}

}
func UpdateFile(context *gin.Context) {

	var basicPapper fileStr
	if err := context.ShouldBindJSON(&basicPapper); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	ghClient := githubclient.NewGitHubClient()
	userInfo, err := utils.GetUserInfo(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	repoName := utils.FormatRepositoryName(userInfo.ID + "_" + basicPapper.Papper)
	Filename := utils.FormatRepositoryName(basicPapper.Path)

	err = ghClient.UpdateFile(repoName, userInfo.Email, userInfo.Name, Filename+"/"+Filename, "initial File", "initial doc")

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "nao foi possivel atualizar arquivo. " + repoName + Filename})
		return
	}

}

type fileStr struct {
	Papper string `json:"papper"`
	Path   string `json:"path"`
}

func CreateRootPapper(C *gin.Context) {
	userInfo, err := utils.GetUserInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user info. " + err.Error()})
		return
	}
	var rp models.RootPapper
	C.ShouldBindJSON(&rp)

	if rp.Name == "" {
		rp.Name = "PapperTrail"
		return
	}

	client, err := userInfo.GetClient(googleOauthConfig)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error file listo. " + err.Error()})
	}

	driveSrv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error file listo. " + err.Error()})
	}
	folder, err := googleClient.CreateFolder(driveSrv, rp.Name, "")
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating root folder.  " + err.Error()})
	}

	rp.Id = folder.Id
	rp.UserID = userInfo.ID

	rp.Save()

	C.JSON(http.StatusOK, gin.H{"message": rp.Name + " folder created successfully."})
}

func CreatePapper(C *gin.Context) {
	var papper models.Papper
	C.ShouldBindJSON(&papper)

	rootPapperInfo, err := utils.GetRootPapperInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting root folder info. " + err.Error()})
		return
	}
	userInfo, err := utils.GetUserInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user info. " + err.Error()})
		return
	}

	client, err := userInfo.GetClient(googleOauthConfig)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error geting google client info. " + err.Error()})
	}

	driveSrv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error geting google driver. " + err.Error()})
	}
	folder, err := googleClient.CreateFolder(driveSrv, papper.Name, rootPapperInfo.Id)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating folder. " + err.Error()})
	}
	papper.Root_papper_id = rootPapperInfo.Id
	papper.Path = rootPapperInfo.Name + "/" + papper.Name
	papper.ID = folder.Id
	papper.Save()
	fileId, err := googleClient.CreateDocxFile(driveSrv, "Chapter 1", folder.Id, "<h1>Chapter 1</h1>")
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating first chapter. " + err.Error()})
	}
	C.JSON(http.StatusOK, fileId)

}

func GetFileList(C *gin.Context) {
	userInfo, err := utils.GetUserInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user info. " + err.Error()})
		return
	}

	client, err := userInfo.GetClient(googleOauthConfig)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error file listo. " + err.Error()})
	}

	driveSrv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error file listo. " + err.Error()})
	}

	fileList, err := googleClient.ListFiles(driveSrv)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error file listo. " + err.Error()})
		return
	}
	C.JSON(http.StatusOK, fileList)

}

func GetDocxIdByProject(C *gin.Context) {

	userInfo, err := utils.GetUserInfo(C)

	baseFolder := C.Query("baseFolder")
	// projectName := C.Query("project")

	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user info. " + err.Error()})
		return
	}

	client, err := userInfo.GetClient(googleOauthConfig)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error file listo. " + err.Error()})
	}

	driveSrv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error file listo. " + err.Error()})
	}

	fileList, err := googleClient.GetDocxFileIDs(driveSrv, baseFolder)
	C.JSON(http.StatusOK, fileList)

}
