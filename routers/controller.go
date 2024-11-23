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

	C.JSON(http.StatusOK, gin.H{"message": world.Name + " folder created successfully.", "worldID": world.Id, "status": "success"})
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

	// {
	//   "type": "service_account",
	//   "project_id": "concise-beanbag-333616",
	//   "private_key_id": "da7151f00621767c1801e89fdbeaf4e522364b6c",
	//   "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCR+sv60EyrL4VD\nTemLztscvLgN0hNdQBJEwgM9Nr11XH7xZMP2SmyQLeaR8Ikc/XW+5Q9TsWI35i6q\n1KC9QXX/AN7VABDz5PMB86qfcX6IseBKbXz6aqlxRf/nh7nXsEKot2Xxsz6ZfQD0\nM+zBXEeDlilyQSeUvsmogJ+LnyeckphMTpqzZtKA3r1XiLuslVMOotIDO1bxMRnn\n88EST0Czv9vrAACyy3Rwd5a2uF9v80QQkHVR9n5YR29GVvYAp/hRQKf4JPsBT7OJ\n95dcqtaxW7WyeAs7ZNz6u2Ouu/57nXnNY8fcE4CwmpNwkM8gUV8z2WYek9VY5W4P\n4gPF101PAgMBAAECggEAB2n97ySiKDWXJpZv7J5aaYi6LlWDj0JgSyaxZGwBzlKe\nzeLIwxr/jYkPQ008oYDL+KCImT8SlnU13I5FBaer9wITzyycL87qeqhl+4gXnZiv\nJAiQhuVg6rRb7WXqzeYRVKFP56krXj9Hi5+RgDaQGUJIo5lkibzw4AJ8V+qC6ARE\nWQUSbILqtoTg+7xiL3uP/Q8ZFb4V5IXsGaYYDPYqJBa52/WxG294SD4AYqtbEmZY\np+wowhb0haCIodjMkt9m7dapB8gckIIOuGjUVuCMCp+e+hkjvprPUnyXuBQuwD67\n4slTXxKYWuCa86N+CUGxl1cg/PChCUoKLntDyOlsaQKBgQDKEkuaDfe+THOV6Lfo\n+pRrwsBULNqJyGFzcwUnWZmmwm2Mh1ZxVnthxgiOoUaNxoMcQaHGPQ2EwNCPbuyX\nUuSfPLUUMCdQWBLsQX1fhgMN7wBlgVW58uCaJFwBCI9pa5vw2Qq4Zqhx3v/ftA7G\nhIyygaiKfwwYxeg+wYpm9lz9eQKBgQC48EJ/UCa59DqVj7SkGORfHYAcU2xT5mmK\nasPIq/bGpbB0LVvfzN3rYtiBg5A4DMaPsfNwQ6fmNCD7k8bwX9qDenwNHRUi8QZi\naouek48QPD+7lP6khW0noFlLI7OJYKYdlZ12BysI3DidVhrENwjx/Z/HC0Y1t1Ln\nqUtTHZqXBwKBgGav1Wt8HaG/CB3uHUdvz2zTkxkzkfrisWMR2FSe2846j6ESRYNj\nB2AwWrjgjBIQByCc2bD75ZrIwTOikuhzX2rsVrjjn5bcqwEUZrncSEEUa4cpqn7M\nRgcO4xJDX12bKavDIAeFY6Q6Rp1PyxJm2Xj9GsEGvwb3y4XYpJSeLbNBAoGAcGVk\npKd7wcwSxs7dxFV0hfIR6CUzUxJX1k3oy07n3fbY9OKUUcHapbIfTyc8QTRSgQZv\noy0bH6dS3FMFtxUqYnnQZs/kBqZhcPK8BBY9/mn/eeuljyugGVM0sZvzA2z/yD8j\nwZW9q9bbeZPZFKM2BoxTzM6nTwIpmq2jH9KAH4UCgYEAn3A9IBzHBO6H3Y2NFjwD\nb3U/Cxq+4S5KMwvdacXaj9rKwTRFlASUMCxS/tlimlUGo5fo/JwMl+HAIFfmv1VK\nhcIjPo05/RJVBHE9CfY+Sg6o6o5mBCvgRp+1BWjy4mKt9LqDl7V7MQc4+6cCGi4P\nZ2KUIPifmb7FR36hMAMVgUg=\n-----END PRIVATE KEY-----\n",
	//   "client_email": "Papertrail@concise-beanbag-333616.iam.gserviceaccount.com",
	//   "client_id": "102028377717608211022",
	//   "auth_uri": "https://accounts.google.com/o/oauth2/auth",
	//   "token_uri": "https://oauth2.googleapis.com/token",
	//   "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
	//   "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/Papertrail%40concise-beanbag-333616.iam.gserviceaccount.com",
	//   "universe_domain": "googleapis.com"
	// }

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
	chapter.Event_Id = nil

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
func CreateEvent(C *gin.Context) {

	var event models.Event

	if err := C.ShouldBindJSON(&event); err != nil {
		C.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := event.Save()
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
	err := gc.Save()
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
