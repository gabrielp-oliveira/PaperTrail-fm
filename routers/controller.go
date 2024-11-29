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
		folder, err := googleClient.CreateFolder(driveSrv, userInfo.Email, userInfo.ID, "")
		if err != nil {
			C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user base folder. " + err.Error()})
			return
		}
		userInfo.Base_folder = folder.Id
		userInfo.UpdateBaseFolder()
	}

	folderId, err := googleClient.CreateFolder(driveSrv, userInfo.Email, world.Name, userInfo.Base_folder)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating world .  " + err.Error()})
		return
	}
	var desc models.Description
	desc.Description_data = world.Description

	world.Id = folderId.Id

	descId := uuid.New().String()
	world.Description = descId
	world.UserID = userInfo.ID
	desc.Resource_id = folderId.Id
	desc.Resource_type = "world"
	world.CreatedAt = time.Now()
	err = world.Save()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving world .  " + err.Error()})
		return
	}
	desc.Id = descId
	err = desc.Save()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving description .  " + err.Error()})
		return
	}

	ssId := uuid.New().String()
	var settings models.Subway_settings
	settings.Id = ssId
	settings.World_id = world.Id
	settings.Chapter_names = false
	settings.Zomm = 1
	settings.X = 1
	settings.Y = 1
	settings.Group_connection_update_chapter = false
	settings.Display_table_chapters = false

	err = settings.Save()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating basic settings for world.  " + err.Error()})
		return
	}

	var paper models.Paper
	paper.Name = "First Paper"
	paper.Description = "First Paper"
	paper.Order = 1
	paper.Path = world.Name + "/" + paper.Name
	paper.World_id = world.Id
	folder, err := googleClient.CreateFolder(driveSrv, userInfo.Email, paper.Name, world.Id)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving paper.  " + err.Error()})
		return
	}
	paper.ID = folder.Id
	err = paper.Create()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving paper.  " + err.Error()})
		return
	}

	var chapter models.Chapter

	chapter.Name = "First Chapter"
	chapter.Description = "First chapter"
	chapter.Order = 1
	chapter.WorldsID = world.Id
	chapter.PaperID = paper.ID
	chapterId, err := googleClient.CreateChapter(driveSrv, chapter.Name, chapter.PaperID, userInfo.Email)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating chapter document.  " + err.Error()})
		return
	}
	chapter.Id = chapterId
	err = chapter.Create()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving chapter.  " + err.Error()})
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

	C.JSON(http.StatusOK, gin.H{"message": world.Name + " folder created successfully.", "world": world, "status": "success"})
}

func DeleteChapter(C *gin.Context) {
	chapterId := C.Query("id")
	var chp models.Chapter
	chp.Id = chapterId
	if chp.Id == "" {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Chapter id is empty."})
		return
	}
	driveSrv := googleClient.GetAppDriveService()
	err := chp.DeleteChapter(driveSrv)

	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "err deletng chapter:" + err.Error()})
		return

	}
	C.JSON(http.StatusOK, gin.H{"message": " chapter Deleted successfully.", "chapterId": chapterId, "status": "success"})
}
func DeletePaper(C *gin.Context) {
	paperId := C.Query("id")
	var pp models.Paper
	pp.ID = paperId
	if pp.ID == "" {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "paper id is empty."})
		return
	}
	driveSrv := googleClient.GetAppDriveService()
	err := pp.Delete(driveSrv)

	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "err deletng paper:" + err.Error()})
		return

	}
	C.JSON(http.StatusOK, gin.H{"message": " paper Deleted successfully.", "paperId": paperId, "status": "success"})
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

	folder, err := googleClient.CreateFolder(driveSrv, userInfo.Email, paper.Name, worldsInfo.Id)

	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating folder. " + err.Error()})
		return
	}

	paper.World_id = worldsInfo.Id
	paper.Path = worldsInfo.Name + "/" + paper.Name
	paper.ID = folder.Id
	paper.Created_at = time.Now()

	err = paper.Create()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// _, err = googleClient.CreateReadmeFile(driveSrv, paper.ID, "# read me content for this paper.")
	// if err != nil {
	// 	C.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating read me. " + err.Error()})
	// 	return
	// }

	var chapter models.Chapter
	chapter.Name = "chapter 1"
	chapter.PaperID = folder.Id
	chapter.WorldsID = worldsInfo.Id
	chapter.CreatedAt = time.Now()
	chapterId, err := googleClient.CreateChapter(driveSrv, chapter.Name, folder.Id, userInfo.Email)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error " + err.Error()})
		return
	}
	chapter.Id = chapterId
	chapter.Order = 1
	chapter.TimelineID = nil
	chapter.Event_Id = nil

	err = chapter.Create()

	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
func CreateEvent(C *gin.Context) {

	var event models.Event

	if err := C.ShouldBindJSON(&event); err != nil {
		C.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var desc models.Description
	Id := uuid.New().String()
	event.Id = Id
	desc.Resource_id = Id

	err := desc.CreateInitialDescription("event")
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	event.Description = desc.Id
	err = event.Save()
	if err != nil {
		C.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	C.JSON(http.StatusOK, event)
}

func GetWorldChapters(C *gin.Context) {

	worlds, err := utils.GetWorldsInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	chapters, err := worlds.GetWorldChapters()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	C.JSON(http.StatusOK, chapters)

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
	// err = chapter.Save()

	descId := uuid.New().String()
	var desc models.Description
	desc.Description_data = chapter.Description
	chapter.Description = descId

	err = chapter.Save()
	desc.Id = descId
	desc.Resource_type = "chapter"
	desc.Resource_id = docId
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = desc.Save()

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

	// var chapter models.Chapter
	// chapter.Id = chapterId
	// err := chapter.Get()
	// if err != nil {
	// 	C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	driveSrv := googleClient.GetAppDriveService()

	// docId := chapter.Id
	file, err := driveSrv.Files.Get(chapterId).Fields("webViewLink").Do()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching document link: " + err.Error()})
		return
	}

	C.JSON(http.StatusOK, gin.H{"documentUrl": file.WebViewLink})
}
func GetChapter(C *gin.Context) {
	chapterId := C.Query("id")
	if chapterId == "" {
		C.JSON(http.StatusBadRequest, gin.H{"error": "Error getting chapter id."})
		return
	}

	var chapter models.Chapter
	chapter.Id = chapterId
	color, err := chapter.Get()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving chapter: " + err.Error()})
		return
	}
	var result models.ChapterDetails
	next, err := chapter.GetChapterByPaperAndOrder(chapter.Order + 1)
	if err == nil {
		result.Next = *next
	}
	prev, err := chapter.GetChapterByPaperAndOrder(chapter.Order - 1)
	if err == nil {
		result.Prev = *prev
	}
	related, err := chapter.GetRelatedChapters()
	if err == nil {
		result.RelatedChapters = related
	}

	result.Chapter = chapter
	result.Color = color
	// Verificar e obter Timeline, se presente
	if chapter.TimelineID != nil {
		var timeline models.Timeline
		timeline.Id = *chapter.TimelineID
		if err := timeline.Get(); err != nil {
			C.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving timeline: " + err.Error()})
			return
		}
		result.Timeline = timeline
	}

	// Verificar e obter StoryLine, se presente
	if chapter.Storyline_id != nil {
		var storyline models.StoryLine
		storyline.Id = *chapter.Storyline_id
		if err := storyline.Get(); err != nil {
			C.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving storyline: " + err.Error()})
			return
		}
		result.StoryLine = storyline
	}

	driveSrv := googleClient.GetAppDriveService()

	file, err := driveSrv.Files.Get(chapter.Id).Fields("webViewLink").Do()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching document link: " + err.Error()})
		return
	}

	result.DocumentUrl = file.WebViewLink
	C.JSON(http.StatusOK, result)
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

	eventId := C.Query("id")

	var evt models.Event
	evt.Id = eventId

	if evt.Id == "" {
		C.JSON(http.StatusInternalServerError, gin.H{"error": "event id is empty."})
		return
	}

	err := evt.Delete()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	C.JSON(http.StatusOK, evt.Id)

}
func UpdateEvent(C *gin.Context) {
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
	C.JSON(http.StatusOK, event)

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

func GetDescription(C *gin.Context) {
	var d models.Description // Suponha que `Element` é a estrutura de seus elementos
	var resource_id = C.Query("resource_id")

	if resource_id == "" {
		C.JSON(http.StatusBadRequest, gin.H{"error": "resource_id empty"})
		return
	}
	d.Resource_id = resource_id
	err := d.GetByResourceId()
	if err != nil {
		C.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	C.JSON(http.StatusOK, d)
}
func UpdateDescription(C *gin.Context) {
	var d models.Description // Suponha que `Element` é a estrutura de seus elementos

	if err := C.ShouldBindJSON(&d); err != nil {
		C.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := d.Update()
	if err != nil {
		C.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	C.JSON(http.StatusOK, d)
}
func updateSettings(C *gin.Context) {
	var ss models.Subway_settings // Suponha que `Element` é a estrutura de seus elementos

	if err := C.ShouldBindJSON(&ss); err != nil {
		C.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := ss.Update()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	C.JSON(http.StatusOK, ss)

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

	// var chapterTimeline models.ChapterTimeline
	// chapterTimeline.Chapter_Id = &chapter.Id
	// chapterTimeline.TimelineID = chapter.TimelineID
	// chapterTimeline.Range = chapter.Range

	// err = chapterTimeline.Update()
	// if err != nil {
	// 	C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
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
	timelineID := uuid.New().String()
	var desc models.Description
	desc.Description_data = tl.Description

	tl.Id = timelineID
	descId := uuid.New().String()
	tl.Description = descId

	err = tl.Save()
	desc.Id = descId
	desc.Resource_type = "timeline"
	desc.Resource_id = timelineID
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = desc.Save()
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
	var desc models.Description
	desc.Description_data = stl.Description
	desc.Resource_id = id
	err = desc.CreateInitialDescription("storyline")
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	stl.Description = desc.Id
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
func UpdateTimelineList(C *gin.Context) {

	var tls []models.Timeline

	if err := C.ShouldBindJSON(&tls); err != nil {
		C.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, str := range tls {
		if err := str.Update(); err != nil {
			C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	C.JSON(http.StatusOK, tls)
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

	worlds, err := utils.GetWorldsInfo(C)
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var stl models.StoryLine
	stl.Id = stlId

	err = stl.Delete()

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

func CreateGroupConnection(C *gin.Context) {

	var gc models.GroupConnection

	if err := C.ShouldBindJSON(&gc); err != nil {
		C.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := uuid.New().String()

	gc.Id = id

	var desc models.Description
	desc.Resource_id = id
	desc.Description_data = gc.Description
	err := desc.CreateInitialDescription("group_connection")
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	gc.Description = desc.Id
	err = gc.Save()
	if err != nil {
		C.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	C.JSON(http.StatusOK, gc)

}

func UpdateConnection(C *gin.Context) {

	var cnn models.Connection

	if err := C.ShouldBindJSON(&cnn); err != nil {
		C.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := cnn.Update()
	if err != nil {
		C.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	C.JSON(http.StatusOK, cnn)

}
func UpdateGroupConnection(C *gin.Context) {

	var gc models.GroupConnection

	if err := C.ShouldBindJSON(&gc); err != nil {
		C.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := gc.Update()
	if err != nil {
		C.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	C.JSON(http.StatusOK, gc)

}
