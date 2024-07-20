package routes

import (
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

	"github.com/gin-gonic/gin"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var googleOauthConfig = googleClient.StartCredentials()

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

	// client, err := userInfo.GetClient(googleOauthConfig)
	// if err != nil {
	// 	C.JSON(http.StatusInternalServerError, gin.H{"error": "Error file listo. " + err.Error()})
	// }

	// driveSrv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	// if err != nil {
	// 	C.JSON(http.StatusInternalServerError, gin.H{"error": "Error file listo. " + err.Error()})
	// }

	driveSrv := googleClient.GetAppDriveService()

	if userInfo.Base_folder == "" {
		folder, err := googleClient.CreateFolder(driveSrv, userInfo.ID, "")
		if err != nil {
			C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user base folder. " + err.Error()})
		}
		userInfo.Base_folder = folder.Id
		userInfo.UpdateBaseFolder()
	}

	folderId, err := googleClient.CreateFolder(driveSrv, rp.Name, userInfo.Base_folder)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating root folder.  " + err.Error()})
	}

	rp.Id = folderId.Id
	rp.UserID = userInfo.ID

	rp.Save()

	C.JSON(http.StatusOK, gin.H{"message": rp.Name + " folder created successfully."})
}

func CreatePapper(C *gin.Context) {
	rootPapperInfo, err := utils.GetRootPapperInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	papper, err := utils.GetPapperInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userInfo, err := utils.GetUserInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user info. " + err.Error()})
		return
	}

	// client, err := userInfo.GetClient(googleOauthConfig)
	// if err != nil {
	// 	C.JSON(http.StatusInternalServerError, gin.H{"error": "Error geting google client info. " + err.Error()})
	// 	return
	// }

	// driveSrv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	driveSrv := googleClient.GetAppDriveService()
	// if err != nil {
	// 	C.JSON(http.StatusInternalServerError, gin.H{"error": "Error geting google driver. " + err.Error()})
	// 	return
	// }
	folder, err := googleClient.CreateFolder(driveSrv, papper.Name, rootPapperInfo.Id)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating folder. " + err.Error()})
		return
	}
	papper.Root_papper_id = rootPapperInfo.Id
	papper.Path = rootPapperInfo.Name + "/" + papper.Name
	papper.ID = folder.Id
	err = papper.Save()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_, err = googleClient.CreateReadmeFile(driveSrv, papper.ID, "# read me content for this papper.")
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating read me. " + err.Error()})
		return
	}

	repoPath := "tempRepositories/" + papper.ID + "/" + "chapter1"
	err = os.MkdirAll(repoPath, os.ModePerm)
	if err != nil {
		gitConfig.RemoveLocalRepo(repoPath)
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}

	docxFilePath := filepath.Join(repoPath, "chapter_1.docx")
	docxContent := "This is the content of Chapter 1."
	err = gitConfig.CreateDocxFile(docxFilePath, docxContent)
	if err != nil {
		gitConfig.RemoveLocalRepo(repoPath)
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}
	repo, err := git.PlainInit(repoPath, false)
	if err != nil {
		gitConfig.RemoveLocalRepo(repoPath)
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}

	worktree, err := repo.Worktree()
	if err != nil {
		gitConfig.RemoveLocalRepo(repoPath)
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}

	_, err = worktree.Add("chapter_1.docx")
	if err != nil {
		gitConfig.RemoveLocalRepo(repoPath)
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
		gitConfig.RemoveLocalRepo(repoPath)
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}

	_, err = repo.CommitObject(commit)
	if err != nil {
		gitConfig.RemoveLocalRepo(repoPath)
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}

	docId, err := gitConfig.UploadDirectoryToDrive(driveSrv, repoPath, papper.ID, "")
	if err != nil {
		gitConfig.RemoveLocalRepo(repoPath)
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}

	var chapter models.Chapter
	chapter.Id = docId
	chapter.Name = "chapter 1"
	chapter.Papper_id = papper.ID
	chapter.Description = "First chapter from " + papper.Name
	chapter.Created_at = time.Now()
	err = chapter.Save()
	if err != nil {
		gitConfig.RemoveLocalRepo(repoPath)
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}
	gitConfig.RemoveLocalRepo(repoPath)
	log.Println("Repository successfully uploaded to Google Drive")

	C.JSON(http.StatusOK, papper.ID)
}

func CreateChapter(C *gin.Context) {

	chapter, err := utils.GetChapterInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	papper, err := utils.GetPapperInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting root folder info. " + err.Error()})
		return
	}
	userInfo, err := utils.GetUserInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user info. " + err.Error()})
		return
	}

	// client, err := userInfo.GetClient(googleOauthConfig)
	// if err != nil {
	// 	C.JSON(http.StatusInternalServerError, gin.H{"error": "Error geting google client info. " + err.Error()})
	// 	return
	// }

	// driveSrv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	// if err != nil {
	// 	C.JSON(http.StatusInternalServerError, gin.H{"error": "Error geting google driver. " + err.Error()})
	// 	return
	// }
	driveSrv := googleClient.GetAppDriveService()

	repoPath := "tempRepositories/" + papper.ID + "/" + chapter.Name
	err = os.MkdirAll(repoPath, os.ModePerm)
	if err != nil {
		gitConfig.RemoveLocalRepo(repoPath)
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}

	docxFilePath := filepath.Join(repoPath, chapter.Name+".docx")
	docxContent := "This is the content of " + chapter.Name
	err = gitConfig.CreateDocxFile(docxFilePath, docxContent)
	if err != nil {
		gitConfig.RemoveLocalRepo(repoPath)
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}
	repo, err := git.PlainInit(repoPath, false)
	if err != nil {
		gitConfig.RemoveLocalRepo(repoPath)
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}

	worktree, err := repo.Worktree()
	if err != nil {
		gitConfig.RemoveLocalRepo(repoPath)
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}

	_, err = worktree.Add(chapter.Name + ".docx")
	if err != nil {
		gitConfig.RemoveLocalRepo(repoPath)
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
		gitConfig.RemoveLocalRepo(repoPath)
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}

	_, err = repo.CommitObject(commit)
	if err != nil {
		gitConfig.RemoveLocalRepo(repoPath)
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}

	docId, err := gitConfig.UploadDirectoryToDrive(driveSrv, repoPath, papper.ID, "")
	if err != nil {
		gitConfig.RemoveLocalRepo(repoPath)
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}
	chapter.Id = docId
	err = chapter.Save()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}
	gitConfig.RemoveLocalRepo(repoPath)
	log.Println("Repository successfully uploaded to Google Drive")
	file, err := driveSrv.Files.Get(docId).Fields("webViewLink").Do()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}
	C.JSON(http.StatusOK, file)

}

func GetRootPapperList(C *gin.Context) {
	userInfo, err := utils.GetUserInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user info. " + err.Error()})
		return
	}
	list, err := userInfo.GetRootPappers()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user root papper. " + err.Error()})
		return
	}
	C.JSON(http.StatusOK, list)

}

func GetChapterList(C *gin.Context) {
	papper, err := utils.GetPapperInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	list, err := papper.GetChapterList()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting chapter. " + err.Error()})
		return
	}
	C.JSON(http.StatusOK, list)

}

func getPapperList(C *gin.Context) {
	rootPapper, err := utils.GetRootPapperInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	list, err := rootPapper.GetPapperList()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting papper. " + err.Error()})
		return
	}
	C.JSON(http.StatusOK, list)

}

func GetChapter(C *gin.Context) {
	userInfo, err := utils.GetUserInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user info. " + err.Error()})
		return
	}

	var chapter models.Chapter
	C.ShouldBindJSON(&chapter)

	url := fmt.Sprintf("https://docs.google.com/document/d/%s/edit?access_token=%s", chapter.Id, userInfo.AccessToken)
	C.JSON(http.StatusOK, url)
}
