package googleClient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var OauthStateString = "randomstatestring"

func StartCredentials() *oauth2.Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(".env load error: %v", err)
	}
	ClientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	if ClientID == "" {
		log.Fatalf("credentials error: %v", err)
	}

	ClientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	if ClientSecret == "" {
		log.Fatalf("credentials error: %v", err)
	}

	return &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		ClientSecret: ClientSecret,
		ClientID:     ClientID,
		Endpoint:     google.Endpoint,
	}
}

func GetGoogleToken(config *oauth2.Config, code string) (*oauth2.Token, error) {
	return config.Exchange(context.Background(), code)
}

func GetAppDriveService() *drive.Service {
	serviceAccountFile := "config/google_client_secret.json"

	// Ler o arquivo JSON da chave da conta de serviço
	b, err := os.ReadFile(serviceAccountFile)
	if err != nil {
		log.Fatalf("Unable to read service account file: %v", err)
	}

	// Autenticar usando a conta de serviço
	config, err := google.JWTConfigFromJSON(b, drive.DriveFileScope)
	if err != nil {
		log.Fatalf("Unable to parse service account file to config: %v", err)
	}

	// Criar cliente HTTP
	client := config.Client(context.Background())

	// Criar serviço do Google Drive
	service, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create drive service: %v", err)
	}

	return service

}

func GetAppDriveServiceWithToken(accessToken string) (*drive.Service, error) {
	ctx := context.Background()
	httpClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: accessToken,
	}))

	srv, err := drive.New(httpClient)
	if err != nil {
		return nil, err
	}
	return srv, nil
}

func VerifyGoogleToken(accessToken string) (*oauth2.Token, error) {
	tokenInfoURL := fmt.Sprintf("https://www.googleapis.com/oauth2/v3/tokeninfo?access_token=%s", accessToken)
	resp, err := http.Get(tokenInfoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid token")
	}

	var tokenInfo oauth2.Token
	err = json.NewDecoder(resp.Body).Decode(&tokenInfo)
	if err != nil {
		return nil, err
	}

	return &tokenInfo, nil
}

func GetGoogleDoc(accessToken, docId string) (*http.Response, error) {
	docURL := fmt.Sprintf("https://docs.googleapis.com/v1/documents/%s", docId)
	req, err := http.NewRequest("GET", docURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	return client.Do(req)
}
