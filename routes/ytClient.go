package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func getClientFromCache(config *oauth2.Config) *http.Client {
	tokenFile := os.ExpandEnv("$HOME/.credentials/youtube-go-quickstart.json")
	f, err := os.Open(tokenFile)
	if err != nil {
		log.Fatalf("Cannot open token file: %v", err)
	}
	defer f.Close()

	tok := &struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		RefreshToken string `json:"refresh_token"`
		Expiry       string `json:"expiry"`
	}{}
	if err := json.NewDecoder(f).Decode(tok); err != nil {
		log.Fatalf("Cannot decode token: %v", err)
	}

	oauthTok := &oauth2.Token{
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
	}
	return config.Client(context.Background(), oauthTok)
}

func getNewYtService() *youtube.Service {
	ctx := context.Background()
	b, err := os.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Cannot read client secret: %v", err)
	}

	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Cannot parse client secret JSON: %v", err)
	}

	client := getClientFromCache(config)
	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Cannot create YouTube service: %v", err)
	}
	return service

}
