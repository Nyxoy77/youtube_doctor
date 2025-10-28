package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nyxoy77/B2C_YouTube_Doctor/constant"
	"github.com/nyxoy77/B2C_YouTube_Doctor/models"
	"github.com/nyxoy77/B2C_YouTube_Doctor/service"
	"github.com/nyxoy77/B2C_YouTube_Doctor/utils"
)

type handler struct {
	service service.DoctorServiceInterface
}

func NewHandler(s service.DoctorServiceInterface) *handler {
	return &handler{service: s}
}

func (h *handler) GetVideos(c *gin.Context) {
	var reqBody models.GetVideoRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("error occurred parsing request body: %v", err)
		utils.WriteError(c, http.StatusBadRequest, err)
		return
	}

	if reqBody.ChannelName == "" {
		utils.WriteError(c, http.StatusBadRequest, constant.MissingChannelError)
		return
	}

	finalResponse, err := h.service.TriggerService(c.Request.Context(), reqBody.ChannelName)
	if err != nil {
		utils.WriteError(c, http.StatusBadRequest, err)
		return
	}
	utils.WriteSuccess(c, http.StatusOK, finalResponse)
}
