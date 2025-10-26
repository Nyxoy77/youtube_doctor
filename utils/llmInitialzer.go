package utils

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/genai"
)

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
