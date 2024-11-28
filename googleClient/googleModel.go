package googleClient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"google.golang.org/api/drive/v3"
)

func ListFiles(service *drive.Service) ([]*drive.File, error) {
	files, err := service.Files.List().PageSize(10).Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		return nil, err
	}

	return files.Files, nil
}

func DeleteFilesInFolder(service *drive.Service, folderId string) error {
	pageToken := ""
	for {
		// Listar arquivos dentro da pasta
		query := fmt.Sprintf("'%s' in parents", folderId)
		files, err := service.Files.List().Q(query).PageToken(pageToken).Do()
		if err != nil {
			return err
		}

		// Excluir os arquivos dentro da pasta
		for _, file := range files.Files {
			fmt.Printf("Deleting file in folder: %s\n", file.Name)
			err := service.Files.Delete(file.Id).Do()
			if err != nil {
				fmt.Printf("Error deleting file %s: %v\n", file.Name, err)
			} else {
				fmt.Printf("Successfully deleted file: %s\n", file.Name)
			}
		}

		// Se houver mais arquivos na pasta, continuar
		if files.NextPageToken == "" {
			break
		}
		pageToken = files.NextPageToken
	}
	return nil
}

// Função para deletar uma pasta vazia
func DeleteFolder(service *drive.Service, folderId string) error {
	// Deletar a pasta
	fmt.Printf("Deleting folder: %s\n", folderId)
	err := service.Files.Delete(folderId).Do()
	if err != nil {
		fmt.Printf("Error deleting folder %s: %v\n", folderId, err)
		return err
	} else {
		fmt.Printf("Successfully deleted folder: %s\n", folderId)
	}
	return nil
}

func DeleteAllFilesAndFolders(service *drive.Service) error {
	pageToken := ""
	for {
		// Fazendo a requisição para listar arquivos e pastas com paginação
		files, _ := service.Files.List().Q("trashed = false").PageToken(pageToken).Do()

		// Excluir arquivos
		for _, file := range files.Files {
			// Se for uma pasta, deletar primeiro seus arquivos internos
			if file.MimeType == "application/vnd.google-apps.folder" {
				DeleteFilesInFolder(service, file.Id) // Deleta os arquivos dentro da pasta
				DeleteFolder(service, file.Id)        // Deleta a pasta
			} else {
				// Se for um arquivo, excluir diretamente
				fmt.Printf("Deleting file: %s\n", file.Name)
				err := service.Files.Delete(file.Id).Do()
				if err != nil {
					fmt.Printf("Error deleting file %s: %v\n", file.Name, err)
					return err
				} else {
					fmt.Printf("Successfully deleted file: %s\n", file.Name)
				}
			}
		}

		// Se houver mais arquivos para listar
		if files.NextPageToken == "" {
			break
		}
		pageToken = files.NextPageToken
	}
	return nil
}

// CreateFolder cria uma pasta no Google Drive e retorna a pasta criada.
func CreateFolder(service *drive.Service, userId, name, parentID string) (*drive.File, error) {
	// Verificar se a pasta já existe
	query := fmt.Sprintf("name='%s' and mimeType='application/vnd.google-apps.folder' and trashed=false", name)
	if parentID != "" {
		query += fmt.Sprintf(" and '%s' in parents", parentID)
	}
	r, err := service.Files.List().Q(query).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to check if folder exists: %v", err)
	}
	if len(r.Files) > 0 {
		return r.Files[0], nil // Retorna a pasta existente em vez de lançar erro
	}

	// Criar a nova pasta
	folder := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
	}

	if parentID != "" {
		folder.Parents = []string{parentID}
	}

	folder, err = service.Files.Create(folder).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to create folder: %v", err)
	}
	return folder, nil
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

func UploadFile(service *drive.Service, name, parentID, filePath string) (string, error) {
	fileMetadata := &drive.File{
		Name:     name,
		MimeType: "application/octet-stream",
	}

	if parentID != "" {
		fileMetadata.Parents = []string{parentID}
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("unable to open file: %v", err)
	}
	defer file.Close()

	_, err = ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("unable to read file content: %v", err)
	}

	createdFile, err := service.Files.Create(fileMetadata).Media(file).Do()
	if err != nil {
		return "", fmt.Errorf("unable to create file: %v", err)
	}

	return createdFile.Id, nil
}

func UploadLocalFile(service *drive.Service, filePath, parentID string) error {
	fileName := filepath.Base(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("unable to open file: %v", err)
	}

	defer file.Close()

	fileMetadata := &drive.File{
		Name:    fileName,
		Parents: []string{parentID},
	}

	_, err = service.Files.Create(fileMetadata).Media(file).Do()
	if err != nil {
		return fmt.Errorf("unable to create file: %v", err)
	}

	return nil
}

func CreateReadmeFile(service *drive.Service, parentID, content string) (*drive.File, error) {
	fileMetadata := &drive.File{
		Name:     "README.md",
		MimeType: "text/markdown",
	}

	if parentID != "" {
		fileMetadata.Parents = []string{parentID}
	}

	readmeFile, err := ioutil.TempFile("", "README.md")
	if err != nil {
		return nil, fmt.Errorf("unable to create temporary file: %v", err)
	}
	defer os.Remove(readmeFile.Name())

	_, err = readmeFile.Write([]byte(content))
	if err != nil {
		return nil, fmt.Errorf("unable to write to temporary file: %v", err)
	}

	readmeFile.Seek(0, 0)
	createdFile, err := service.Files.Create(fileMetadata).Media(readmeFile).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to create README.md file: %v", err)
	}

	return createdFile, nil
}

func CreateChapter(service *drive.Service, name, PaperId, userEmail string) (string, error) {
	folder := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{PaperId},
	}
	createdFolder, err := service.Files.Create(folder).Do()
	if err != nil {
		// Caso o arquivo não seja encontrado, retornamos um erro mais específico
		return "", fmt.Errorf("invalid PaperId (folder not found): %w", err)
	}

	// Criar o documento diretamente dentro da pasta do Paper (livro)
	fileMetadata := &drive.File{
		Name:     name,
		Parents:  []string{createdFolder.Id},
		MimeType: "application/vnd.google-apps.document",
	}

	// Criação do documento
	doc, err := service.Files.Create(fileMetadata).Do()
	if err != nil {
		return "", fmt.Errorf("error creating document: %w", err)
	}

	// Definir permissões para o usuário no documento criado
	permission := &drive.Permission{
		Type:         "user",
		Role:         "writer",
		EmailAddress: userEmail,
	}

	// Aplicar as permissões ao documento
	_, err = service.Permissions.Create(doc.Id, permission).Do()
	if err != nil {
		return "", fmt.Errorf("error setting document permissions: %w", err)
	}

	// Retorna o ID do documento criado
	return doc.Id, nil
}
