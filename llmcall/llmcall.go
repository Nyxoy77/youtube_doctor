package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/nyxoy77/B2C_YouTube_Doctor/models"
	"google.golang.org/genai"
)

func AnalyzeVideoTitles(titles []string) (*models.Response, error) {
	ctx := context.Background()
	client, err := initialIzeGeminiClient(ctx)
	if err != nil {
		return nil, err
	}

	prompt := `Analyze these YouTube video titles and provide improvements. 
    Return the response in this exact JSON format for each title:
    {
        "llmResponse": [
            {
                "prevTitle": "original title",
                "newTitle": "improved title",
                "reason": "explanation for the improvement"
            }
        ]
    }
    
    Only return valid JSON, no additional text.
    Here are the titles to analyze:%s`

	finalPrompt := fmt.Sprintf(prompt, strings.Join(titles, "\n"))
	response, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(finalPrompt), nil)
	if err != nil {
		return nil, fmt.Errorf("error from llm call: %w", err)
	}
	var resp *models.Response
	text := strings.TrimSpace(response.Text())
	start := strings.IndexByte(text, '{')
	end := strings.LastIndexByte(text, '}')
	if start == -1 || end == -1 || start > end {
		return nil, fmt.Errorf("error from LLM call, could not find any valid josn object: %w", err)
	}
	jsonBlob := text[start : end+1]
	fmt.Println(jsonBlob)
	if err := json.Unmarshal([]byte(jsonBlob), &resp); err != nil {
		return nil, fmt.Errorf("errror from llm call failed to unmarshal : %w", err)
	}

	return resp, nil

}

func initialIzeGeminiClient(ctx context.Context) (*genai.Client, error) {
	config := &genai.ClientConfig{
		APIKey: os.Getenv("API_KEY"),
	}
	client, err := genai.NewClient(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error setting up the client %w", err)
	}
	return client, nil
}
