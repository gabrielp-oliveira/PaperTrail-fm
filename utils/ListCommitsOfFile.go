package utils

import (
	"errors"
	"sort"
	"time"

	"github.com/google/go-github/github"
)

type Author struct {
	Date  string `json:"date"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Commit struct {
	Author  Author `json:"author"`
	Message string `json:"message"`
}

type ReducedCommit struct {
	SHA    string `json:"sha"`
	Commit Commit `json:"commit"`
}

func FilterCommits(commits []*github.RepositoryCommit) []ReducedCommit {
	var filteredCommits []ReducedCommit
	for _, commit := range commits {
		if commit.SHA == nil || commit.Commit == nil || commit.Commit.Author == nil {
			continue
		}
		// Converte time.Time para string
		dateString := commit.Commit.Author.GetDate().Format(time.RFC3339)

		filteredCommit := ReducedCommit{
			SHA: *commit.SHA,
			Commit: Commit{
				Author: Author{
					Date:  dateString,
					Name:  commit.Commit.Author.GetName(),
					Email: commit.Commit.Author.GetEmail(),
				},
				Message: commit.Commit.GetMessage(),
			},
		}
		filteredCommits = append(filteredCommits, filteredCommit)
	}
	return filteredCommits
}
func FindPreviousSHA(commits []ReducedCommit, sha string) (string, error) {
	// Ordenar os commits pela data
	sort.Slice(commits, func(i, j int) bool {
		dateI, errI := time.Parse(time.RFC3339, commits[i].Commit.Author.Date)
		dateJ, errJ := time.Parse(time.RFC3339, commits[j].Commit.Author.Date)
		if errI != nil || errJ != nil {
			return false
		}
		return dateI.Before(dateJ)
	})

	// Encontrar o commit com o SHA fornecido e retornar o SHA anterior
	for i := 1; i < len(commits); i++ {
		if commits[i].SHA == sha {
			return commits[i-1].SHA, nil
		}
	}

	return "", errors.New("SHA não encontrado ou é o primeiro commit")
}

type FileChanges struct {
	Sha       string
	Filename  string
	Additions *int `json:"additions,omitempty"`
	Deletions *int `json:"deletions,omitempty"`
	Changes   *int `json:"changes,omitempty"`
	Status    string
	Patch     string
	Diff      *[]string `json:"diff,omitempty"`
	Docx      string    `json:"docx,omitempty"`
}
