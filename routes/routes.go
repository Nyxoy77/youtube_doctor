package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	llm "github.com/nyxoy77/B2C_YouTube_Doctor/llmcall"
	"github.com/nyxoy77/B2C_YouTube_Doctor/models"
)

func GetVideos(c *gin.Context) {
	var reqBody models.GetVideoRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("error occurred parsing request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	log.Println(reqBody.ChannelName)
	if reqBody.ChannelName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing channel name",
		})
		return
	}

	channelID, err := getChannelID(reqBody.ChannelName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	titles, err := fetchVideos(channelID)
	if err != nil {
		fmt.Println("AAAAA")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	finalResponse, err := llm.AnalyzeVideoTitles(titles)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error in llm call": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"response": finalResponse,
	})
}

func getChannelID(channelName string) (string, error) {
	service := getNewYtService()

	// First try with ForUsername
	call := service.Channels.List([]string{"id"}).ForUsername(channelName)
	response, err := call.Do()
	if err == nil && len(response.Items) > 0 {
		return response.Items[0].Id, nil
	}

	// If username search fails, try searching by channel title
	searchCall := service.Search.List([]string{"id"}).
		Type("channel").
		Q(channelName).
		MaxResults(1)

	searchResponse, err := searchCall.Do()
	if err != nil {
		return "", fmt.Errorf("error searching for channel: %w", err)
	}

	if len(searchResponse.Items) == 0 {
		return "", fmt.Errorf("no channel found with name: %s", channelName)
	}

	return searchResponse.Items[0].Id.ChannelId, nil
}

func fetchVideos(channelID string) ([]string, error) {
	service := getNewYtService()
	searchCall := service.Search.List([]string{"id", "snippet"}).ChannelId(channelID).MaxResults(5)

	searchResponse, err := searchCall.Do()
	fmt.Println(err)
	if err != nil {
		return nil, fmt.Errorf("error occurred in the YouTube search call: %w", err)
	}

	var titles []string

	for _, item := range searchResponse.Items {
		titles = append(titles, item.Snippet.Title)
	}
	return titles, nil
}
