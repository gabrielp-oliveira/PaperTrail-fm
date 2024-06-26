package githubclient

import (
	"context"
	"log"
	"os"

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
	gh.repo = repo
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
func (gh *GitHubClient) GetFileContents(path string) ([]byte, *github.RepositoryContent, error) {
	fileContent, resp, _, err := gh.client.Repositories.GetContents(gh.ctx, gh.owner, gh.repo, path, nil)
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
func (gh *GitHubClient) UpdateFile(repo string, userEmail string, userName string, path, message, content, sha string) error {
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
	_, fileContent, err := gh.GetFileContents(path)
	if err != nil {
		log.Printf("Arquivo não encontrado, criando novo: %v", err)
		return gh.CreateFile(repo, userEmail, userName, path, message, content)
	}
	sha := fileContent.GetSHA()
	return gh.UpdateFile(repo, userEmail, userName, path, message, content, sha)
}

// ListCommitsOfFile lista os commits que modificaram um arquivo específico no repositório.
func (gh *GitHubClient) ListCommitsOfFile(path string) ([]*github.RepositoryCommit, error) {
	opts := &github.CommitsListOptions{
		Path: path,
	}
	commits, _, err := gh.client.Repositories.ListCommits(gh.ctx, gh.owner, gh.repo, opts)
	if err != nil {
		log.Printf("Erro ao listar commits: %v", err)
		return nil, err
	}
	return commits, nil
}
