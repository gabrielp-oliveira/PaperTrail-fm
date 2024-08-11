package routes

import (
	"database/sql"
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
	"github.com/google/uuid"
)

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

func CreateWorld(C *gin.Context) {
	userInfo, err := utils.GetUserInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user info. " + err.Error()})
		return
	}
	var rp models.World
	C.ShouldBindJSON(&rp)

	if rp.Name == "" {
		rp.Name = "PapperTrail"
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
			return
		}
		userInfo.Base_folder = folder.Id
		userInfo.UpdateBaseFolder()
	}

	folderId, err := googleClient.CreateFolder(driveSrv, rp.Name, userInfo.Base_folder)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating world .  " + err.Error()})
		return
	}

	rp.Id = folderId.Id
	rp.UserID = userInfo.ID
	rp.CreatedAt = time.Now()
	rp.Save()

	C.JSON(http.StatusOK, gin.H{"message": rp.Name + " folder created successfully."})
}

func CreatePapper(C *gin.Context) {
	worldsInfo, err := utils.GetWorldsInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	papper, err := utils.GetPapperInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// userInfo, err := utils.GetUserInfo(C)
	// if err != nil {
	// 	C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user info. " + err.Error()})
	// 	return
	// }

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
	folder, err := googleClient.CreateFolder(driveSrv, papper.Name, worldsInfo.Id)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating folder. " + err.Error()})
		return
	}
	papper.World_id = worldsInfo.Id
	papper.Path = worldsInfo.Name + "/" + papper.Name
	papper.ID = folder.Id
	papper.Created_at = time.Now()
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

	// repoPath := "tempRepositories/" + papper.ID + "/" + "chapter1"
	// err = os.MkdirAll(repoPath, os.ModePerm)
	// if err != nil {
	// 	gitConfig.RemoveLocalRepo(repoPath)
	// 	C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
	// 	return
	// }
	chapterFolder, err := googleClient.CreateFolder(driveSrv, "Chapter 1", papper.ID)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}
	docxContent := "This is the content of Chapter 1."
	chapterId, err := googleClient.CreateDocxFile(driveSrv, "Chapter 1", chapterFolder.Id, docxContent)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}

	var chapter models.Chapter
	chapter.Id = chapterId
	chapter.Name = "chapter 1"
	chapter.PapperID = papper.ID
	chapter.WorldsID = worldsInfo.Id
	chapter.Description = "First chapter from " + papper.Name
	chapter.CreatedAt = time.Now()

	eventID := sql.NullString{String: ""}
	timelineID := sql.NullString{String: ""}

	chapter.TimelineID = timelineID
	chapter.EventID = eventID

	err = chapter.Save()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}
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
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting world  info. " + err.Error()})
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
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}

	docxFilePath := filepath.Join(repoPath, chapter.Name+".docx")
	docxContent := "This is the content of " + chapter.Name
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

	_, err = worktree.Add(chapter.Name + ".docx")
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

	docId, err := googleClient.CreateChapter(driveSrv, chapter.Name, chapter.PapperID, userInfo.Email)

	if err != nil {
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

func GetWorldsList(C *gin.Context) {
	userInfo, err := utils.GetUserInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user info. " + err.Error()})
		return
	}
	list, err := userInfo.GetWorldss()
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
	userInfo, err := utils.GetUserInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	driveSrv := googleClient.GetAppDriveService()

	list, err := papper.GetChapterList(driveSrv, userInfo.AccessToken)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting chapter. " + err.Error()})
		return
	}
	C.JSON(http.StatusOK, list)

}

func getPapperList(C *gin.Context) {
	worlds, err := utils.GetWorldsInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	list, err := worlds.GetPapperList()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting papper. " + err.Error()})
		return
	}
	C.JSON(http.StatusOK, list)

}

func GetChapter(C *gin.Context) {
	ChapterInfo, err := utils.GetChapterInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user info. " + err.Error()})
		return
	}
	driveSrv := googleClient.GetAppDriveService()

	revisions, err := driveSrv.Revisions.List(ChapterInfo.Id).Do()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve revisions: " + err.Error()})

	}

	// url := fmt.Sprintf("https://docs.google.com/document/d/%s/edit?access_token=%s", ChapterInfo.Id, userInfo.AccessToken)
	C.JSON(http.StatusOK, revisions)
}

func GetWorldData(C *gin.Context) {
	worlds, err := utils.GetWorldsInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	environment, err := worlds.GetWorldData()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	C.JSON(http.StatusOK, environment)

}

func InsertEvent(C *gin.Context) {
	worlds, err := utils.GetWorldsInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	event, err := utils.GetEventInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id := uuid.New().String()

	event.Id = id
	event.Save()

	environment, err := worlds.GetWorldData()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	C.JSON(http.StatusOK, environment)

}
func RemoveEvent(C *gin.Context) {
	worlds, err := utils.GetWorldsInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	event, err := utils.GetEventInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if event.Id == "" {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "event id is empty."})
		return
	}

	err = event.Delete()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	environment, err := worlds.GetWorldData()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	C.JSON(http.StatusOK, environment)

}
func UpdateEvent(C *gin.Context) {
	worlds, err := utils.GetWorldsInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	event, err := utils.GetEventInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = event.Update()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	environment, err := worlds.GetWorldData()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	C.JSON(http.StatusOK, environment)

}
func InsertConnection(C *gin.Context) {
	worlds, err := utils.GetWorldsInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	cnn, err := utils.GetConnectionInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id := uuid.New().String()
	cnn.Id = id
	err = cnn.Save()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	environment, err := worlds.GetWorldData()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	C.JSON(http.StatusOK, environment)

}
func RemoveConnection(C *gin.Context) {
	worlds, err := utils.GetWorldsInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	cnn, err := utils.GetConnectionInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if cnn.Id == "" {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "connection Id is empty."})
		return
	}
	err = cnn.Delete()
	if cnn.Id == "" {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	environment, err := worlds.GetWorldData()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	C.JSON(http.StatusOK, environment)

}
