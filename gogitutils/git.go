package gogitutils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Repository struct {
	repo *git.Repository
}

func CreateNewRepository(repoDir string) (*Repository, error) {
	// Criar o diretório se ele não existir
	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		os.MkdirAll(repoDir, os.ModePerm)
	}

	repo, err := git.PlainInit(repoDir, false)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar novo repositório: %v", err)
	}

	return &Repository{repo: repo}, nil
}

func (r *Repository) AddFile(rootDir, fileName string) error {
	wt, err := r.repo.Worktree()
	if err != nil {
		return fmt.Errorf("falha ao obter Worktree: %v", err)
	}

	// Caminho completo do arquivo na raiz da aplicação
	srcFilePath := filepath.Join(rootDir, fileName)

	// Caminho para onde o arquivo será copiado dentro do repositório
	destFilePath := filepath.Join(wt.Filesystem.Root(), fileName)

	// Copiar o arquivo para o diretório do repositório
	err = copyFile(srcFilePath, destFilePath)
	if err != nil {
		return fmt.Errorf("falha ao copiar arquivo para o repositório: %v", err)
	}

	// Adicionar o arquivo ao índice
	_, err = wt.Add(fileName)
	if err != nil {
		return fmt.Errorf("falha ao adicionar arquivo %q ao índice: %v", fileName, err)
	}

	return nil
}

func (r *Repository) Commit(message, authorName, authorEmail string) error {
	wt, err := r.repo.Worktree()
	if err != nil {
		return fmt.Errorf("falha ao obter Worktree: %v", err)
	}

	_, err = wt.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  authorName,
			Email: authorEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("falha ao realizar commit: %v", err)
	}

	return nil
}

func (r *Repository) ListCommits() error {
	ref, err := r.repo.Head()
	if err != nil {
		return fmt.Errorf("falha ao obter head: %v", err)
	}

	iter, err := r.repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return fmt.Errorf("falha ao obter log: %v", err)
	}

	err = iter.ForEach(func(c *object.Commit) error {
		fmt.Println(c)
		return nil
	})
	if err != nil {
		return fmt.Errorf("falha ao iterar commits: %v", err)
	}

	return nil
}

func (r *Repository) SetUser(name, email string) error {
	cfg, err := r.repo.Config()
	if err != nil {
		return fmt.Errorf("falha ao obter configuração do repositório: %v", err)
	}

	// if cfg.User == nil {
	// 	cfg.User = &config.User{}
	// }

	cfg.User.Name = name
	cfg.User.Email = email

	return r.repo.Storer.SetConfig(cfg)
}

// Função auxiliar para copiar arquivos
func copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("falha ao abrir arquivo fonte: %v", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("falha ao criar arquivo de destino: %v", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("falha ao copiar conteúdo do arquivo: %v", err)
	}

	return nil
}
