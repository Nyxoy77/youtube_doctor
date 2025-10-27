package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"google.golang.org/genai"
)

func WriteError(c *gin.Context, status int, message any) {
	c.JSON(status, gin.H{
		"error": message,
	})
}

func WriteSuccess(c *gin.Context, status int, message any) {
	c.JSON(status, gin.H{
		"response": message,
	})
}

func InitialIzeGeminiClient(ctx context.Context) (*genai.Client, error) {
	config := &genai.ClientConfig{
		APIKey: os.Getenv("API_KEY"),
	}
	client, err := genai.NewClient(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error setting up the client %w", err)
	}
	return client, nil
}

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

func GlobalLoggerInstance() *logrus.Logger {
	logger := &logrus.Logger{
		Out: os.Stderr,
		Formatter: &logrus.TextFormatter{
			FullTimestamp: true,
		},
	}
	return logger
}
