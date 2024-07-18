package models

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"PaperTrail-fm.com/db"
	"golang.org/x/oauth2"
)

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email" binding:"required"`
	Name         string    `json:"name"`
	Password     string    `json:"password"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refresh_token"`
	TokenExpiry  time.Time `json:"token_expiry"`
	Base_folder  string    `json:"base_folder"`
}

func (u User) GetClient(config *oauth2.Config) (*http.Client, error) {
	var token oauth2.Token

	// Recupere o token do banco de dados
	err := db.DB.QueryRow("SELECT accessToken, refresh_token, token_expiry FROM users WHERE email = $1", u.Email).Scan(
		&token.AccessToken, &token.RefreshToken, &token.Expiry)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from database: %v", err)
	}

	// Verifique se o token está expirado e, se necessário, use o refresh token para obter um novo token
	if time.Now().After(token.Expiry) {
		tokenSource := config.TokenSource(context.Background(), &token)
		newToken, err := tokenSource.Token()
		if err != nil {
			return nil, fmt.Errorf("unable to refresh token: %v", err)
		}
		// Atualize o token no banco de dados
		_, err = db.DB.Exec("UPDATE users SET accessToken = $1, refresh_token = $2, token_expiry = $3 WHERE email = $4",
			newToken.AccessToken, newToken.RefreshToken, newToken.Expiry, u.Email)
		if err != nil {
			return nil, fmt.Errorf("unable to update token in database: %v", err)
		}
		token = *newToken
	}

	client := config.Client(context.Background(), &token)

	return client, nil
}

func (u User) UpdateToken() error {
	updateQuery := "UPDATE users SET accessToken = $1, refresh_token = $2, token_expiry = $3 WHERE email = $4"
	_, err := db.DB.Exec(updateQuery, u.AccessToken, u.RefreshToken, u.TokenExpiry, u.Email)
	if err != nil {
		return errors.New("Error updating token. " + err.Error())
	}
	return nil
}
func (u User) UpdateBaseFolder() error {
	updateQuery := "UPDATE users SET base_folder = $1 WHERE email = $2"
	_, err := db.DB.Exec(updateQuery, u.Base_folder, u.Email)
	if err != nil {
		return errors.New("Error updating base folder. " + err.Error())
	}
	return nil
}

func (u User) SetToken() error {
	insertQuery := "INSERT INTO users(email, password, created_at, id, name, accessToken, refresh_token, token_expiry) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"

	_, err := db.DB.Exec(insertQuery,
		u.Email, "", time.Now(), u.ID, u.Name, u.AccessToken, u.RefreshToken, u.TokenExpiry)
	if err != nil {
		return errors.New("Unable to save token in database: " + err.Error())

	}
	return nil
}
