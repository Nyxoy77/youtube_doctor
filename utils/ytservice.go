package utils

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func GetNewYtService() *youtube.Service {
	ctx := context.Background()

	b, err := os.ReadFile("client_secret.json")
	if err != nil {
		log.Printf("Unable to read client_secret.json: %v", err)
		return nil
	}

	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Printf("Unable to parse client secret: %v", err)
		return nil
	}

	tokenFile := os.ExpandEnv("$HOME/.credentials/youtube-go-quickstart.json")
	f, err := os.Open(tokenFile)
	if err != nil {
		log.Printf("Cannot open token file: %v", err)
		return nil
	}
	defer f.Close()

	tok := &oauth2.Token{}
	if err := json.NewDecoder(f).Decode(tok); err != nil {
		log.Printf("Cannot decode token: %v", err)
		return nil
	}

	client := config.Client(ctx, tok)

	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Printf("Error creating YouTube service: %v", err)
		return nil
	}

	return service
}
