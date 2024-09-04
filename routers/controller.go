package routes

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"PaperTrail-fm.com/googleClient"
	"PaperTrail-fm.com/models"
	"PaperTrail-fm.com/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetCommitDiff(context *gin.Context) {

	// repo := context.Query("repo")
	// sha := context.Query("sha")
	// path := context.Query("path")
	// ghClient := githubclient.NewGitHubClient()

	// userInfo, err := utils.GetUserInfo(context)
	// if err != nil {
	// 	context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// repoName := utils.FormatRepositoryName(userInfo.ID + "_" + repo)

	// commits, err := ghClient.GetCommitDiff(repoName, sha, path)
	// if err != nil {
	// 	context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// context.JSON(http.StatusOK, commits)

}

func CreateWorld(C *gin.Context) {
	userInfo, err := utils.GetUserInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user info. " + err.Error()})
		return
	}
	var world models.World
	C.ShouldBindJSON(&world)

	if world.Name == "" {
		world.Name = "PapperTrail"
	}

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

	folderId, err := googleClient.CreateFolder(driveSrv, world.Name, userInfo.Base_folder)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating world .  " + err.Error()})
		return
	}

	world.Id = folderId.Id
	world.UserID = userInfo.ID
	world.CreatedAt = time.Now()
	err = world.Save()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving world .  " + err.Error()})
		return
	}

	var strline models.StoryLine
	_, err = strline.CreateBasicStoryLines(world.Id)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving story lines.  " + err.Error()})
		return
	}
	var timeline models.Timeline
	_, err = timeline.CreateBasicTimelines(world.Id)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving time lines.  " + err.Error()})
		return
	}
	C.JSON(http.StatusOK, gin.H{"message": world.Name + " folder created successfully."})
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

	userInfo, err := utils.GetUserInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user info. " + err.Error()})
		return
	}

	driveSrv := googleClient.GetAppDriveService()
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

	var chapter models.Chapter
	chapter.Name = "chapter 1"
	chapter.PapperID = papper.ID
	chapter.WorldsID = worldsInfo.Id
	chapter.Description = "First chapter from " + papper.Name
	chapter.CreatedAt = time.Now()
	chapterId, err := googleClient.CreateChapter(driveSrv, chapter.Name, chapter.PapperID, userInfo.Email)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}
	chapter.Id = chapterId
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
	type papperChapter struct {
		models.Papper
		Chapter []models.Chapter `json:"chapter"`
	}

	response := papperChapter{
		Papper: papper,
	}
	response.Chapter = append(response.Chapter, chapter)
	C.JSON(http.StatusOK, response)
}

func CreateChapter(C *gin.Context) {

	chapter, err := utils.GetChapterInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting chapter info: " + err.Error()})
		return
	}

	driveSrv := googleClient.GetAppDriveService()

	userInfo, err := utils.GetUserInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user info: " + err.Error()})
		return
	}

	// Cria o documento no Google Drive
	docId, err := googleClient.CreateChapter(driveSrv, chapter.Name, chapter.PapperID, userInfo.Email)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating document: " + err.Error()})
		return
	}

	// Atualiza o ID do capítulo e salva
	chapter.Id = docId
	err = chapter.Save()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving chapter: " + err.Error()})
		return
	}

	file, err := driveSrv.Files.Get(docId).Fields("webViewLink").Do()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching document link: " + err.Error()})
		return
	}

	C.JSON(http.StatusOK, gin.H{"documentUrl": file.WebViewLink})
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

func GetChapterUrl(C *gin.Context) {
	chapterId := C.Query("id")
	if chapterId == "" {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting chapter id."})
		return
	}

	var chapter models.Chapter
	chapter.Id = chapterId
	err := chapter.Get()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	driveSrv := googleClient.GetAppDriveService()

	docId := chapter.Id
	// chapter.Get()
	file, err := driveSrv.Files.Get(docId).Fields("webViewLink").Do()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching document link: " + err.Error()})
		return
	}

	C.JSON(http.StatusOK, gin.H{"documentUrl": file.WebViewLink})
}
func GetChapter(C *gin.Context) {
	chapterId := C.Query("id")
	if chapterId == "" {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting chapter id."})
		return
	}

	var chapter models.Chapter
	chapter.Id = chapterId

	err := chapter.Get()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting chapter info: " + err.Error()})
		return
	}

	C.JSON(http.StatusOK, chapter)
}
func GetPapper(C *gin.Context) {
	papperId := C.Query("id")
	if papperId == "" {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting chapter id."})
		return
	}

	var papper models.Papper
	papper.ID = papperId

	err := papper.Get()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting papper info: " + err.Error()})
		return
	}

	C.JSON(http.StatusOK, papper)
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
func UpdatePapper(C *gin.Context) {
	papper, err := utils.GetPapperInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = papper.Update()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	C.JSON(http.StatusOK, papper)

}
func UpdateChapter(C *gin.Context) {
	chapter, err := utils.GetChapterInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	currentTime := time.Now()
	chapter.LastUpdate = &currentTime
	err = chapter.UpdateChapter()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	C.JSON(http.StatusOK, chapter)

}
func CreateConnection(C *gin.Context) {
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
