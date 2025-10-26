package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nyxoy77/B2C_YouTube_Doctor/models"
	"github.com/nyxoy77/B2C_YouTube_Doctor/service"
)

type Handler struct {
	service service.DoctorServiceInterface
}

func NewHandler(s service.DoctorServiceInterface) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetVideos(c *gin.Context) {
	var reqBody models.GetVideoRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("error occurred parsing request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if reqBody.ChannelName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing channel name",
		})
		return
	}

	finalResponse, err := h.service.TriggerService(c.Request.Context(), reqBody.ChannelName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error from service": err,
		})
		return
	}
	c.JSON(200, gin.H{
		"response": finalResponse,
	})
}
