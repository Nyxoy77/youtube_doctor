package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nyxoy77/B2C_YouTube_Doctor/constant"
	"github.com/nyxoy77/B2C_YouTube_Doctor/models"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/youtube/v3"
	"google.golang.org/genai"
)

type DoctorClientInterface interface {
	FetchTitles(ctx context.Context, channelName string) ([]string, error)
	LlmCall(ctx context.Context, titles []string) (*models.Response, error)
}

type DoctorClient struct {
	ytService   *youtube.Service
	genAiClient *genai.Client
	logger      *logrus.Logger
}

func NewDoctorClient(ytService *youtube.Service, genAiClient *genai.Client, logger *logrus.Logger) DoctorClientInterface {
	return &DoctorClient{
		ytService:   ytService,
		genAiClient: genAiClient,
		logger:      logger,
	}
}

func (c *DoctorClient) FetchTitles(ctx context.Context, channelName string) ([]string, error) {
	searchCall := c.ytService.Search.List([]string{"id"}).
		Context(ctx).
		Type(constant.Channel).
		Q(channelName).
		MaxResults(1)

	searchResponse, err := searchCall.Do()
	if err != nil {
		return []string{}, fmt.Errorf("error searching for channel: %w", err)
	}

	if len(searchResponse.Items) == 0 {
		return []string{}, fmt.Errorf("no channel found with name: %s", channelName)
	}

	channelID := searchResponse.Items[0].Id.ChannelId

	searchCall = c.ytService.Search.List([]string{"id", "snippet"}).ChannelId(channelID).MaxResults(5)

	searchResponse, err = searchCall.Do()

	if err != nil {
		return nil, fmt.Errorf("error occurred in the YouTube search call: %w", err)
	}

	var titles []string

	for _, item := range searchResponse.Items {
		titles = append(titles, item.Snippet.Title)
	}
	return titles, nil
}

func (c *DoctorClient) LlmCall(ctx context.Context, titles []string) (*models.Response, error) {

	finalPrompt := fmt.Sprintf(constant.Prompt, strings.Join(titles, "\n"))

	response, err := c.genAiClient.Models.GenerateContent(ctx, constant.Gemini_Flash, genai.Text(finalPrompt), nil)
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

	if err := json.Unmarshal([]byte(jsonBlob), &resp); err != nil {
		return nil, fmt.Errorf("errror from llm call failed to unmarshal : %w", err)
	}

	return resp, nil
}
