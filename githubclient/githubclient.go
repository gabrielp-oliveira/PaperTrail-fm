package githubclient

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"PaperTrail-fm.com/utils"
	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

// GitHubClient é uma estrutura que mantém o cliente GitHub e o contexto.
type GitHubClient struct {
	client *github.Client
	ctx    context.Context
	owner  string
	repo   string
}

// NewGitHubClient cria uma nova instância de GitHubClient autenticada com o token fornecido.
func NewGitHubClient() *GitHubClient {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	token := os.Getenv("GITHUB_USER_KEY")
	if token == "" {
		log.Fatalf("Token de usuário do GitHub não encontrado no ambiente")
	}

	owner := os.Getenv("GITHUB_USER_NAME")
	if owner == "" {
		log.Fatalf("Nome de usuário do GitHub não encontrado no ambiente")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return &GitHubClient{
		client: client,
		ctx:    ctx,
		owner:  owner,
	}
}

// SetRepo define o repositório que será utilizado por operações subsequentes.
func (gh *GitHubClient) SetRepo(repo string) {
	repo = repo
}

// CreateRepo cria um novo repositório no GitHub.
func (gh *GitHubClient) CreateRepo(name, description string, private bool) (*github.Repository, error) {
	repo := &github.Repository{
		Name:        github.String(name),
		Description: github.String(description),
		Private:     github.Bool(private),
	}
	createdRepo, resp, err := gh.client.Repositories.Create(gh.ctx, "", repo)
	if err != nil {
		log.Printf("Erro ao criar repositório: %v, Resposta: %v", err, resp)
		return nil, err
	}
	return createdRepo, nil
}

// GetFileContents obtém o conteúdo binário de um arquivo específico de um repositório.
func (gh *GitHubClient) GetFileContents(path string, repo string) ([]byte, *github.RepositoryContent, error) {
	fileContent, resp, _, err := gh.client.Repositories.GetContents(gh.ctx, gh.owner, repo, path, nil)
	if err != nil {
		log.Printf("Erro ao obter conteúdo do arquivo: %v, Resposta: %v", err, resp)
		return nil, nil, err
	}
	content, err := fileContent.GetContent()
	if err != nil {
		log.Printf("Erro ao obter conteúdo do arquivo: %v", err)
		return nil, nil, err
	}
	return []byte(content), fileContent, nil
}

// CreateFile cria um arquivo em um repositório.
func (gh *GitHubClient) CreateFile(repo string, userEmail string, userName string, path, message, content string) error {
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(message),
		Content: []byte(content),
		Committer: &github.CommitAuthor{
			Name:  github.String(userName),
			Email: github.String(userEmail),
		},
		Author: &github.CommitAuthor{
			Name:  github.String(userName),
			Email: github.String(userEmail),
		},
	}
	_, resp, err := gh.client.Repositories.CreateFile(gh.ctx, gh.owner, repo, path, opts)
	if err != nil {
		log.Printf("Erro ao criar o arquivo: %v, Resposta: %v", err, resp)
		return err
	}
	return nil
}

// UpdateFile atualiza um arquivo em um repositório.
func (gh *GitHubClient) UpdateFile(repo string, userEmail string, userName string, path, message, content string) error {

	_, fileContent, err := gh.GetFileContents(path, repo)
	if err != nil {
		return err
	}
	sha := fileContent.GetSHA()

	opts := &github.RepositoryContentFileOptions{
		Message: github.String(message),
		Content: []byte(content),
		SHA:     github.String(sha),
		Committer: &github.CommitAuthor{
			Name:  github.String(userName),
			Email: github.String(userEmail),
		},
		Author: &github.CommitAuthor{
			Name:  github.String(userName),
			Email: github.String(userEmail),
		},
	}
	_, resp, err := gh.client.Repositories.UpdateFile(gh.ctx, gh.owner, repo, path, opts)
	if err != nil {
		log.Printf("Erro ao atualizar o arquivo: %v, Resposta: %v", err, resp)
		return err
	}
	return nil
}

// DeleteRepo deleta um repositório no GitHub.
func (gh *GitHubClient) DeleteRepo(repoName string) error {
	resp, err := gh.client.Repositories.Delete(gh.ctx, gh.owner, repoName)
	if err != nil {
		log.Printf("Erro ao deletar repositório: %v, Resposta: %v", err, resp)
		return err
	}
	return nil
}

// CreateOrUpdateFile cria ou atualiza um arquivo em um repositório.
func (gh *GitHubClient) CreateOrUpdateFile(repo string, userEmail string, userName string, path string, message, content string) error {
	_, _, err := gh.GetFileContents(path, repo)
	if err != nil {
		log.Printf("Arquivo não encontrado, criando novo: %v", err)
		return gh.CreateFile(repo, userEmail, userName, path, message, content)
	}
	return gh.UpdateFile(repo, userEmail, userName, path, message, content)
}
func (gh *GitHubClient) createFile(repo string, userEmail string, userName string, path string) error {
	_, _, err := gh.GetFileContents(path, repo)
	if err != nil {
		log.Printf("Arquivo não encontrado, criando novo: %v", err)
		return gh.CreateFile(repo, userEmail, userName, path, "initial File", "")
	}
	return nil
}

// ListCommitsOfFile lista os commits que modificaram um arquivo específico no repositório.
func (gh *GitHubClient) ListCommitsOfFile(repo string, path string) ([]utils.ReducedCommit, error) {
	opts := &github.CommitsListOptions{
		Path: path,
	}
	commits, _, err := gh.client.Repositories.ListCommits(gh.ctx, gh.owner, repo, opts)
	if err != nil {
		return nil, err
	}
	filteredCommits := utils.FilterCommits(commits)

	return filteredCommits, nil
}

func (gh *GitHubClient) GetCommitDiff(repo string, sha string, path string) ([]*utils.FileChanges, error) {

	commit, _, err := gh.client.Repositories.GetCommit(gh.ctx, gh.owner, repo, sha)
	if err != nil {
		log.Printf("Error getting commit: %v", err)
		return nil, err
	}

	var files []*utils.FileChanges
	var file utils.FileChanges
	for _, f := range commit.Files {
		if *f.Filename == path {
			file.Filename = *f.Filename
			file.Sha = *f.SHA
			file.Additions = f.Additions
			file.Deletions = f.Deletions
			file.Changes = f.Changes
			file.Status = *f.Status
		}

		if filepath.Ext(*f.Filename) == ".docx" {

			currentTempFilePath, _, err := gh.GetFileFromCommitTemp(repo, sha, path)
			if err != nil {
				return nil, fmt.Errorf("error getting current file content: %v", err)
			}

			previousTempFilePath, err := gh.GetPreviousFileTemp(repo, sha, path)
			if err != nil {
				return nil, fmt.Errorf("error getting previous file content: %v", err)
			}

			diff, err := GetDocxDiff(currentTempFilePath, previousTempFilePath)
			file.Diff = diff
			if err != nil {
				return nil, fmt.Errorf(" %v", err)
			}
		} else {
			if f.Patch != nil {
				file.Patch = *f.Patch
			}
		}
		files = append(files, &file)

	}

	if len(files) == 0 {
		return nil, errors.New("file not found in the specified commit")
	}
	return files, nil
}

func (gh *GitHubClient) GetFileFromCommit(repo string, sha string, path string) ([]byte, string, error) {
	commit, _, err := gh.client.Repositories.GetCommit(gh.ctx, gh.owner, repo, sha)
	if err != nil {
		return nil, "", fmt.Errorf("error getting commit: %v", err)
	}
	var fileSHA string
	for _, file := range commit.Files {
		if *file.Filename == path {
			fileSHA = *file.SHA
			break
		}
	}
	if fileSHA == "" {
		return nil, "", fmt.Errorf("file not found in the specified commit")
	}
	blob, _, err := gh.client.Git.GetBlob(gh.ctx, gh.owner, repo, fileSHA)
	if err != nil {
		return nil, "", fmt.Errorf("error getting blob: %v", err)
	}
	content, err := base64.StdEncoding.DecodeString(*blob.Content)
	if err != nil {
		return nil, "", fmt.Errorf("error decoding content: %v", err)
	}
	return content, filepath.Ext(path), nil
}

func (gh *GitHubClient) GetPreviousFileTemp(repo string, sha string, path string) (string, error) {

	commitList, err := gh.ListCommitsOfFile(repo, path)
	if err != nil {
		return "", err
	}

	lastSha, err := utils.FindPreviousSHA(commitList, sha)
	if err != nil {
		return "", err
	}

	filePath, _, err := gh.GetFileFromCommitTemp(repo, lastSha, path)

	if err != nil {
		return "", err
	}
	return filePath, nil
}

func (gh *GitHubClient) GetFileFromCommitTemp(repo string, sha string, path string) (string, string, error) {
	commit, _, err := gh.client.Repositories.GetCommit(gh.ctx, gh.owner, repo, sha)
	if err != nil {
		return "", "", fmt.Errorf("error getting commit: %v", err)
	}

	var fileSHA string
	for _, file := range commit.Files {
		if *file.Filename == path {
			fileSHA = *file.SHA
			break
		}
	}

	if fileSHA == "" {
		return "", "", fmt.Errorf("file not found in the specified commit")
	}

	blob, _, err := gh.client.Git.GetBlob(gh.ctx, gh.owner, repo, fileSHA)
	if err != nil {
		return "", "", fmt.Errorf("error getting blob: %v", err)
	}

	content, err := base64.StdEncoding.DecodeString(*blob.Content)
	if err != nil {
		return "", "", fmt.Errorf("error decoding content: %v", err)
	}

	// Cria um arquivo temporário
	tmpFile, err := ioutil.TempFile("", "document-*.docx")
	if err != nil {
		return "", "", fmt.Errorf("error creating temporary file: %v", err)
	}
	defer tmpFile.Close()

	// Escreve o conteúdo decodificado no arquivo temporário
	if _, err := tmpFile.Write(content); err != nil {
		return "", "", fmt.Errorf("error writing to temporary file: %v", err)
	}

	// Retorna o caminho do arquivo temporário e sua extensão
	return tmpFile.Name(), filepath.Ext(path), nil
}
