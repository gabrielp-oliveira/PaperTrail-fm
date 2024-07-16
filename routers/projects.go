package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"PaperTrail-fm.com/gitConfig"
	"PaperTrail-fm.com/githubclient"
	"PaperTrail-fm.com/googleClient"
	"PaperTrail-fm.com/models"
	"PaperTrail-fm.com/utils"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	"github.com/gin-gonic/gin"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

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
	folderId, err := googleClient.CreateFolder(driveSrv, rp.Name, "")
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating root folder.  " + err.Error()})
	}

	rp.Id = folderId.Id
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
		return
	}

	driveSrv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error geting google driver. " + err.Error()})
		return
	}
	folder, err := googleClient.CreateFolder(driveSrv, papper.Name, rootPapperInfo.Id)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating folder. " + err.Error()})
		return
	}
	papper.Root_papper_id = rootPapperInfo.Id
	papper.Path = rootPapperInfo.Name + "/" + papper.Name
	papper.ID = folder.Id
	papper.Save()
	_, err = googleClient.CreateReadmeFile(driveSrv, papper.ID, "# read me content for this papper.")
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating read me. " + err.Error()})
		return
	}
	// chapterId, err := googleClient.CreateFolder(driveSrv, "chapter1", papper.ID)
	// if err != nil {
	// 	C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating chapter folder. " + err.Error()})
	// 	return
	// }

	repoPath := "tempRepositories/" + papper.ID + "/" + "chapter1"
	err = os.MkdirAll(repoPath, os.ModePerm)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}

	docxFilePath := filepath.Join(repoPath, "chapter_1.docx")
	docxContent := "This is the content of Chapter 1."
	err = gitConfig.CreateDocxFile(docxFilePath, docxContent)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}
	repo, err := git.PlainInit(repoPath, false)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}

	worktree, err := repo.Worktree()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}

	_, err = worktree.Add("chapter_1.docx")
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}

	commit, err := worktree.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  userInfo.Name,
			Email: userInfo.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}

	_, err = repo.CommitObject(commit)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}

	err = gitConfig.UploadDirectoryToDrive(driveSrv, repoPath, papper.ID)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}
	log.Println("Repository successfully uploaded to Google Drive")

	// fileId, err := googleClient.CreateDocxFile(driveSrv, "Chapter 1", chapterId, "<h1>Chapter 1/h1>")
	// if err != nil {
	// 	C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating chapter docs. " + err.Error()})
	// 	return
	// }
	// repoPath, err := gitConfig.CreateLocalRepo(papper.Path + "/Chapter1", "Chapter1")
	// if err != nil {
	// 	C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating chapter docs. " + err.Error()})
	// 	gitConfig.RemoveLocalRepo(repoPath)
	// 	return
	// }
	// defer func() {
	// 	if err != nil {
	// 		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating chapter docs. " + err.Error()})
	// 		gitConfig.RemoveLocalRepo(repoPath)
	// 		return
	// 	}
	// }()

	// err = gitConfig.UploadRepo(driveSrv, repoPath, folderId)
	// if err != nil {
	// 	gitConfig.RemoveLocalRepo(repoPath)
	// 	log.Fatalf("Unable to upload repo: %v", err)
	// 	return
	// }

	// _, err = gitConfig.UploadFile(driveSrv, "Chapter1", chapterId, papper.Path+"/chapter1")
	// if err != nil {
	// 	C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating first chapter. " + err.Error()})
	// 	return
	// }
	gitConfig.RemoveLocalRepo("repoPath")
	C.JSON(http.StatusOK, "fileId")

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
