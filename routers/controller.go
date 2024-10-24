package routes

import (
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
		world.Name = "PaperTrail"
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

func CreatePaper(C *gin.Context) {
	worldsInfo, err := utils.GetWorldsInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	paper, err := utils.GetPaperInfo(C)
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
	folder, err := googleClient.CreateFolder(driveSrv, paper.Name, worldsInfo.Id)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating folder. " + err.Error()})
		return
	}

	paper.World_id = worldsInfo.Id
	paper.Path = worldsInfo.Name + "/" + paper.Name
	paper.ID = folder.Id
	paper.Created_at = time.Now()
	err = paper.Save()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_, err = googleClient.CreateReadmeFile(driveSrv, paper.ID, "# read me content for this paper.")
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating read me. " + err.Error()})
		return
	}

	var chapter models.Chapter
	chapter.Name = "chapter 1"
	chapter.PaperID = paper.ID
	chapter.WorldsID = worldsInfo.Id
	chapter.Description = "First chapter from " + paper.Name
	chapter.CreatedAt = time.Now()
	chapterId, err := googleClient.CreateChapter(driveSrv, chapter.Name, chapter.PaperID, userInfo.Email)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}
	chapter.Id = chapterId

	chapter.TimelineID = nil
	chapter.EventID = nil

	err = chapter.Save()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}
	log.Println("Repository successfully uploaded to Google Drive")
	type PaperChapter struct {
		models.Paper
		Chapter []models.Chapter `json:"chapter"`
	}

	response := PaperChapter{
		Paper: paper,
	}
	response.Chapter = append(response.Chapter, chapter)
	C.JSON(http.StatusOK, response)
}

func CreateChapter(C *gin.Context) {

	chapter, err := utils.GetContextInfo[models.Chapter](C, "chapter")

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
	docId, err := googleClient.CreateChapter(driveSrv, chapter.Name, chapter.PaperID, userInfo.Email)
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

	// file, err := driveSrv.Files.Get(docId).Fields("webViewLink").Do()
	// if err != nil {
	// 	C.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching document link: " + err.Error()})
	// 	return
	// }

	C.JSON(http.StatusOK, chapter)
}

func GetWorldsList(C *gin.Context) {
	userInfo, err := utils.GetUserInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user info. " + err.Error()})
		return
	}
	list, err := userInfo.GetWorldss()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user root paper. " + err.Error()})
		return
	}
	C.JSON(http.StatusOK, list)

}

func GetChapterList(C *gin.Context) {
	paper, err := utils.GetPaperInfo(C)
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

	list, err := paper.GetChapterList(driveSrv, userInfo.AccessToken)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting chapter. " + err.Error()})
		return
	}
	C.JSON(http.StatusOK, list)

}

func getPaperList(C *gin.Context) {
	worlds, err := utils.GetWorldsInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	list, err := worlds.GetPaperList()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting paper. " + err.Error()})
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
func GetPaper(C *gin.Context) {
	PaperId := C.Query("id")
	if PaperId == "" {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting chapter id."})
		return
	}

	var paper models.Paper
	paper.ID = PaperId

	err := paper.Get()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting paper info: " + err.Error()})
		return
	}

	C.JSON(http.StatusOK, paper)
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
	event, err := utils.GetContextInfo[models.Event](C, "event")

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

	event, err := utils.GetContextInfo[models.Event](C, "event")
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

	event, err := utils.GetContextInfo[models.Event](C, "event")
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
func UpdatePaper(C *gin.Context) {
	paper, err := utils.GetPaperInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = paper.Update()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	C.JSON(http.StatusOK, paper)

}
func UpdatePaperList(C *gin.Context) {
	var pprs []models.Paper // Suponha que `Element` é a estrutura de seus elementos

	if err := C.ShouldBindJSON(&pprs); err != nil {
		C.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, pp := range pprs {
		// Lógica de atualização para cada elemento
		if err := pp.Update(); err != nil {
			C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	C.JSON(http.StatusOK, true)

}
func UpdateChapter(C *gin.Context) {

	chapter, err := utils.GetContextInfo[models.Chapter](C, "chapter")
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	currentTime := time.Now()

	chapter.LastUpdate = &currentTime

	// chapterTL.chapter.TimelineID = &chapterTL.TimelineID
	// chapterTL.chapter.Storyline_id = &chapterTL.Storyline_id

	err = chapter.UpdateChapter()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var chapterTimeline models.ChapterTimeline
	chapterTimeline.Chapter_Id = &chapter.Id
	chapterTimeline.TimelineID = chapter.TimelineID
	chapterTimeline.Range = chapter.Range

	err = chapterTimeline.Update()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	C.JSON(http.StatusOK, chapter)

}
func UpdateChapterList(C *gin.Context) {

	var chpts []models.Chapter

	if err := C.ShouldBindJSON(&chpts); err != nil {
		C.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, chp := range chpts {
		// Lógica de atualização para cada elemento
		currentTime := time.Now()
		chp.LastUpdate = &currentTime
		if err := chp.UpdateChapter(); err != nil {
			C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// chapterTL.chapter.TimelineID = &chapterTL.TimelineID
	// chapterTL.chapter.Storyline_id = &chapterTL.Storyline_id

	C.JSON(http.StatusOK, chpts)

}
func CreateConnection(C *gin.Context) {

	cnn, err := utils.GetContextInfo[models.Connection](C, "connection")
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

	C.JSON(http.StatusOK, cnn)

}
func CreateTimeline(C *gin.Context) {

	tl, err := utils.GetContextInfo[models.Timeline](C, "timeline")
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id := uuid.New().String()
	tl.Id = id
	err = tl.Save()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	C.JSON(http.StatusOK, tl)
}
func CreateStoryline(C *gin.Context) {

	stl, err := utils.GetContextInfo[models.StoryLine](C, "storyline")
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id := uuid.New().String()
	stl.Id = id
	stl.Created_at = time.Now()

	err = stl.Save()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	C.JSON(http.StatusOK, stl)
}
func UpdateTimeline(C *gin.Context) {

	tl, err := utils.GetContextInfo[models.Timeline](C, "timeline")
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = tl.Update()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	C.JSON(http.StatusOK, tl)
}

func UpdateStoryline(C *gin.Context) {

	stl, err := utils.GetContextInfo[models.StoryLine](C, "storyline")
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = stl.Update()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	C.JSON(http.StatusOK, stl)
}
func UpdateStorylineList(C *gin.Context) {

	var strl []models.StoryLine

	if err := C.ShouldBindJSON(&strl); err != nil {
		C.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, str := range strl {
		if err := str.Update(); err != nil {
			C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	C.JSON(http.StatusOK, strl)
}
func DeleteTimeline(C *gin.Context) {
	tlId := C.Query("id")

	var tl models.Timeline
	tl.Id = tlId

	err := tl.Delete()

	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	C.JSON(http.StatusOK, tl)
}
func DeleteStoryline(C *gin.Context) {
	stlId := C.Query("id")

	var stl models.StoryLine
	stl.Id = stlId

	err := stl.Delete()

	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	C.JSON(http.StatusOK, stl)
}
func RemoveConnection(C *gin.Context) {

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

	C.JSON(http.StatusOK, err)

}
