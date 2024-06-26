package utils

import (
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

type FileChanges struct {
	Sha       string
	Filename  string
	Additions *int `json:"additions,omitempty"`
	Deletions *int `json:"deletions,omitempty"`
	Changes   *int `json:"changes,omitempty"`
	Status    string
	Patch     string
}
