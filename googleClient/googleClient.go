package googleClient

import (
	"context"
	"log"
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
