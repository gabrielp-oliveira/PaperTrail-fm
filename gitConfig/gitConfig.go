package gitConfig

import (
	"fmt"
	"io/ioutil"
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

// UploadDirectoryToDrive faz upload recursivo de um diretório local para o Google Drive
func UploadDirectoryToDrive(service *drive.Service, localPath, parentFolderID string) error {
	// Obter informações sobre o diretório local
	fileInfo, err := os.Stat(localPath)
	if err != nil {
		return fmt.Errorf("unable to read directory info: %v", err)
	}

	// Verificar se é um diretório válido
	if !fileInfo.IsDir() {
		return fmt.Errorf("%s is not a valid directory", localPath)
	}

	// Criar uma pasta no Google Drive com o nome do diretório local
	folder := &drive.File{
		Name:     fileInfo.Name(),
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentFolderID},
	}
	createdFolder, err := service.Files.Create(folder).Do()
	if err != nil {
		return fmt.Errorf("unable to create folder in Google Drive: %v", err)
	}
	folderID := createdFolder.Id

	// Listar arquivos no diretório local
	files, err := ioutil.ReadDir(localPath)
	if err != nil {
		return fmt.Errorf("unable to read directory: %v", err)
	}

	// Iterar sobre os arquivos e subdiretórios
	for _, file := range files {
		filePath := filepath.Join(localPath, file.Name())

		if file.IsDir() {
			// Recursivamente fazer upload de subdiretórios
			err := UploadDirectoryToDrive(service, filePath, folderID)
			if err != nil {
				log.Printf("unable to upload directory '%s': %v", filePath, err)
				// Continuar com o próximo arquivo ou diretório
				continue
			}
		} else {
			// Fazer upload de arquivos individuais
			fileMetadata := &drive.File{
				Name:    file.Name(),
				Parents: []string{folderID},
			}
			fileContent, err := os.Open(filePath)
			if err != nil {
				log.Printf("unable to open file '%s': %v", filePath, err)
				// Continuar com o próximo arquivo ou diretório
				continue
			}

			_, err = service.Files.Create(fileMetadata).Media(fileContent).Do()
			fileContent.Close()
			if err != nil {
				log.Printf("unable to upload file '%s': %v", filePath, err)
				// Continuar com o próximo arquivo ou diretório
				continue
			}
		}
	}

	return nil
}
func CreateDocxFile(filePath string, content string) error {
	doc := document.New()

	// Add content to the document
	doc.AddParagraph().AddRun().AddText(content)

	// Save the document to file
	err := doc.SaveToFile(filePath)
	if err != nil {
		return fmt.Errorf("unable to save docx file: %v", err)
	}

	return nil
}
