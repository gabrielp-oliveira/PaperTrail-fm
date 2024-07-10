package googleClient

import (
	"bytes"
	"fmt"

	"google.golang.org/api/drive/v3"
)

func ListFiles(service *drive.Service) ([]*drive.File, error) {
	files, err := service.Files.List().PageSize(10).Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		return nil, err
	}

	return files.Files, nil
}

func CreateFolder(service *drive.Service, name, parentID string) (*drive.File, error) {
	// Verificar se a pasta já existe
	query := fmt.Sprintf("name='%s' and mimeType='application/vnd.google-apps.folder' and trashed=false", name)
	existingFolders, err := service.Files.List().Q(query).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to check existing folders: %v", err)
	}

	if len(existingFolders.Files) > 0 {
		return nil, fmt.Errorf("folder '%s' already exists, please chose another name.", name)
	}

	// Criar a nova pasta
	folder := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
	}

	if parentID != "" {
		folder.Parents = []string{parentID}
	}

	createdFolder, err := service.Files.Create(folder).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to create folder: %v", err)
	}

	return createdFolder, nil
}

func GetDocxFileIDs(service *drive.Service, projectFolderID string) ([]string, error) {
	var docxFileIDs []string
	query := fmt.Sprintf("'%s' in parents and mimeType='application/vnd.google-apps.folder'", projectFolderID)

	// Obter subpastas (capítulos) do projeto
	subFolders, err := service.Files.List().Q(query).Do()
	if err != nil {
		return nil, err
	}

	// Iterar sobre cada subpasta (capítulo)
	for _, subFolder := range subFolders.Files {
		capID := subFolder.Id
		query = fmt.Sprintf("'%s' in parents and mimeType='application/vnd.openxmlformats-officedocument.wordprocessingml.document'", capID)

		// Obter arquivos .docx na subpasta
		files, err := service.Files.List().Q(query).Do()
		if err != nil {
			return nil, err
		}

		// Adicionar IDs dos arquivos .docx à lista
		for _, file := range files.Files {
			docxFileIDs = append(docxFileIDs, file.Id)
		}
	}

	return docxFileIDs, nil
}

func CreateDocxFile(service *drive.Service, name, parentID, content string) (string, error) {
	// Define the file metadata
	fileMetadata := &drive.File{
		Name:     name,
		MimeType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	}

	// Set the parent folder if provided
	if parentID != "" {
		fileMetadata.Parents = []string{parentID}
	}

	// Create the file content as a byte buffer
	fileContent := []byte(content)
	fileReader := bytes.NewReader(fileContent)

	// Create the file on Google Drive
	createdFile, err := service.Files.Create(fileMetadata).Media(fileReader).Do()
	if err != nil {
		return "", fmt.Errorf("unable to create file: %v", err)
	}

	return createdFile.Id, nil
}
