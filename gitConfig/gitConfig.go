package gitConfig

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"PaperTrail-fm.com/googleClient"
	"baliance.com/gooxml/document"
	"google.golang.org/api/drive/v3"
)

// GitClient é uma estrutura que mantém o caminho do repositório local.
type GitClient struct {
	repoPath string
}

// NewGitClient cria uma nova instância de GitClient com o caminho do repositório fornecido.
func NewGitClient(repoPath string) *GitClient {
	return &GitClient{repoPath: repoPath}
}

// CreateLocalRepo cria um novo repositório Git local.
func CreateLocalRepo(repoName, fileName string) (string, error) {
	repoPath := filepath.Join("tempRepositories", repoName)
	err := os.MkdirAll(repoPath, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("could not create directory: %v", err)
	}

	docxFilePath := filepath.Join(repoPath, fileName+".docx")
	docxContent := "This is the content of" + fileName
	err = CreateDocxFile(docxFilePath, docxContent)
	if err != nil {
		log.Fatalf("Unable to create .docx file: %v", err)
	}

	return repoPath, nil
}

// UploadGitConfigFiles faz o upload dos arquivos de configuração do Git para o Google Drive.
func UploadGitConfigFiles(service *drive.Service, repoPath, parentID string) error {
	configFiles := []string{"config", "HEAD", "description"}

	for _, fileName := range configFiles {
		filePath := filepath.Join(repoPath, ".git", fileName)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			continue // Skip if file does not exist
		}

		fileID, err := googleClient.UploadFile(service, fileName, parentID, filePath)
		if err != nil {
			return fmt.Errorf("could not upload file %s: %v", fileName, err)
		}

		log.Printf("Uploaded file %s with ID %s", fileName, fileID)
	}

	return nil
}

func RemoveLocalRepo(repoPath string) error {
	err := os.RemoveAll(repoPath)
	if err != nil {
		return fmt.Errorf("unable to remove local repo: %v", err)
	}
	return nil
}

func UploadDirectoryToDrive(service *drive.Service, userEmail, localPath, parentFolderID, docId string) (string, error) {
	fileInfo, err := os.Stat(localPath)
	if err != nil {
		return "", fmt.Errorf("unable to read directory info: %v", err)
	}

	if !fileInfo.IsDir() {
		return "", fmt.Errorf("%s is not a valid directory", localPath)
	}

	folder := &drive.File{
		Name:     fileInfo.Name(),
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentFolderID},
	}
	createdFolder, err := service.Files.Create(folder).Do()
	if err != nil {
		return "", fmt.Errorf("unable to create folder in Google Drive: %v", err)
	}
	folderID := createdFolder.Id

	files, err := os.ReadDir(localPath)
	if err != nil {
		return "", fmt.Errorf("unable to read directory: %v", err)
	}

	for _, file := range files {
		filePath := filepath.Join(localPath, file.Name())

		if file.IsDir() {
			dcId, err := UploadDirectoryToDrive(service, userEmail, filePath, folderID, docId)
			if err != nil {
				log.Printf("unable to upload directory '%s': %v", filePath, err)
				continue
			}
			docId = dcId
		} else {
			fileMetadata := &drive.File{
				Name:    file.Name(),
				Parents: []string{folderID},
			}

			// Determine the mime type based on file extension
			mimeType := "application/octet-stream"
			if filepath.Ext(filePath) == ".docx" {
				mimeType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
			}
			fileMetadata.MimeType = mimeType
			fileContent, err := os.Open(filePath)
			if err != nil {
				log.Printf("unable to open file '%s': %v", filePath, err)
				continue
			}

			doc, err := service.Files.Create(fileMetadata).Media(fileContent).Do()
			if filepath.Ext(filePath) == ".docx" {
				permission := &drive.Permission{
					Type:         "user",
					Role:         "writer",
					EmailAddress: userEmail,
				}

				_, err = service.Permissions.Create(doc.Id, permission).Do()
				if err != nil {
					log.Fatalf("Unable to create permission: %v", err)
				}
			}
			fileContent.Close()
			if err != nil {
				log.Printf("unable to upload file '%s': %v", filePath, err)
				continue
			}

			// If the file is a .docx, convert it to Google Docs format
			if filepath.Ext(file.Name()) == ".docx" {
				convertFileMetadata := &drive.File{
					MimeType: "application/vnd.google-apps.document",
				}
				_, err = service.Files.Update(doc.Id, convertFileMetadata).Do()
				if err != nil {
					log.Printf("unable to convert file '%s' to Google Docs format: %v", filePath, err)
					continue
				}
				docId = doc.Id
			}
		}
	}

	return docId, nil
}
func CreateDocxFile(filePath string, content string) error {
	doc := document.New()

	para := doc.AddParagraph()
	run := para.AddRun()
	run.AddText(content)

	err := doc.SaveToFile(filePath)
	if err != nil {
		return fmt.Errorf("unable to save docx file: %v", err)
	}

	return nil
}
