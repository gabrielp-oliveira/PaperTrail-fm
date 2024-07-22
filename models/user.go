package models

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"PaperTrail-fm.com/db"
	"PaperTrail-fm.com/googleClient"
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

var googleOauthConfig = googleClient.StartCredentials()

func (u *User) UpdateToken() error {
	updateQuery := "UPDATE users SET accessToken = $1, refresh_token = $2, token_expiry = $3 WHERE email = $4"
	_, err := db.DB.Exec(updateQuery, u.AccessToken, u.RefreshToken, u.TokenExpiry, u.Email)
	if err != nil {
		return errors.New("Error updating token. " + err.Error())
	}
	return nil
}

func (u *User) UpdateBaseFolder() error {
	updateQuery := "UPDATE users SET base_folder = $1 WHERE email = $2"
	_, err := db.DB.Exec(updateQuery, u.Base_folder, u.Email)
	if err != nil {
		return errors.New("Error updating base folder. " + err.Error())
	}
	return nil
}

func (u *User) GetRootPappers() ([]RootPapper, error) {
	query := "SELECT id, name FROM rootpappers WHERE user_id = $1"
	rows, err := db.DB.Query(query, u.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []RootPapper
	for rows.Next() {
		var rp RootPapper
		if err := rows.Scan(&rp.Id, &rp.Name); err != nil {
			return nil, err
		}
		list = append(list, rp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenExpiry  time.Time `json:"token_expiry"`
}

func (u *User) updateDatabase() error {
	updateQuery := "UPDATE users SET name = $1, password = $2, accessToken = $3, refresh_token = $4, token_expiry = $5 WHERE id = $6"
	_, err := db.DB.Exec(updateQuery, u.Name, u.Password, u.AccessToken, u.RefreshToken, u.TokenExpiry, u.ID)
	if err != nil {
		return errors.New("Error updating user. " + err.Error())
	}
	return nil
}

func (u *User) UpdateOAuthToken() (*oauth2.Token, error) {
	config := googleOauthConfig

	token := &oauth2.Token{
		AccessToken:  u.AccessToken,
		RefreshToken: u.RefreshToken,
		Expiry:       u.TokenExpiry,
	}

	tokenSource := config.TokenSource(context.Background(), token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}

	u.AccessToken = newToken.AccessToken
	u.RefreshToken = newToken.RefreshToken
	u.TokenExpiry = newToken.Expiry

	err = u.updateDatabase()
	if err != nil {
		return nil, err
	}

	return newToken, nil
}

func (u *User) GetClient(config *oauth2.Config) (*http.Client, error) {
	var token oauth2.Token

	err := db.DB.QueryRow("SELECT accessToken, refresh_token, token_expiry FROM users WHERE email = $1", u.Email).Scan(
		&token.AccessToken, &token.RefreshToken, &token.Expiry)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from database: %v", err)
	}

	if time.Now().After(token.Expiry) {
		tokenSource := config.TokenSource(context.Background(), &token)
		newToken, err := tokenSource.Token()
		if err != nil {
			return nil, fmt.Errorf("unable to refresh token: %v", err)
		}

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
