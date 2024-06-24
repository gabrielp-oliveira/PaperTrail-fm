package models

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"PaperTrail-fm.com/s3Instance"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

func UploadFileToS3(filePath, fileName string, destination string) error {

	sess := s3Instance.S3Creds.Session

	key := fileName
	if destination != "" {
		key = fmt.Sprintf("%s/%s", destination, fileName)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("falha ao abrir arquivo %q: %v", filePath, err)
	}
	defer file.Close()

	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3Instance.S3Creds.BucketName),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		return fmt.Errorf("falha ao enviar arquivo para o S3: %v", err)
	}

	fmt.Printf("Arquivo %q enviado com sucesso para o bucket %q\n", filePath, s3Instance.S3Creds.BucketName)
	return nil
}

func DownloadFileFromS3(fileName, downloadPath string) error {

	sess := s3Instance.S3Creds.Session

	svc := s3.New(sess)

	// Obter o objeto do S3
	output, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3Instance.S3Creds.BucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return fmt.Errorf("falha ao obter o objeto do S3: %v", err)
	}
	defer output.Body.Close()

	// Criar o diretório onde o arquivo será baixado, se necessário
	if err := os.MkdirAll(filepath.Dir(downloadPath), 0755); err != nil {
		return fmt.Errorf("falha ao criar diretório para %q: %v", downloadPath, err)
	}

	// Criar o arquivo local para armazenar o conteúdo baixado
	file, err := os.Create(downloadPath)
	if err != nil {
		return fmt.Errorf("falha ao criar arquivo local %q: %v", downloadPath, err)
	}
	defer file.Close()

	// Copiar o conteúdo do objeto do S3 para o arquivo local
	if _, err := io.Copy(file, output.Body); err != nil {
		return fmt.Errorf("falha ao copiar conteúdo do objeto do S3 para %q: %v", downloadPath, err)
	}

	return nil
}

func Upload(context *gin.Context) {
	// filePath := "./file1.txt" // Caminho do arquivo no root da aplicação
	// fileName := "file1.txt"

	// err := UploadFileToS3(filePath, fileName)
	// if err != nil {
	// 	log.Fatalf("Erro: %v", err)
	// }
}
func Download(context *gin.Context) {
	fileName := "file1.txt"
	downloadPath := "./temporaryFiles/userId/file1" // Caminho onde o arquivo será salvo localmente

	err := DownloadFileFromS3(fileName, downloadPath)
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}
}

type RequestBody struct {
	UserID string `json:"userId" binding:"required"`
}

func CreateEmptyFolder(context *gin.Context) {

	var requestBody RequestBody

	// Bind JSON ao requestBody e validar
	if err := context.ShouldBindJSON(&requestBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId := requestBody.UserID

	UploadFileToS3("new Papper.docx", "new Papper.docx", userId+"/node-1")
}
