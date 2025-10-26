package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nyxoy77/B2C_YouTube_Doctor/client"
	"github.com/nyxoy77/B2C_YouTube_Doctor/routes"
	"github.com/nyxoy77/B2C_YouTube_Doctor/service"
	"github.com/nyxoy77/B2C_YouTube_Doctor/utils"
)

func main() {
	router := gin.Default()

	ctx := context.Background()
	err := godotenv.Load()

	if err != nil {
		fmt.Printf("error loading the env : %v", err)
		return
	}

	ytService := utils.GetNewYtService()
	genAiClient, err := utils.InitialIzeGeminiClient(ctx)
	if err != nil {
		fmt.Println("error initlaizing the geminiClient", err)
		return
	}
	client := client.NewDoctorClient(ytService, genAiClient)
	svc := service.NewDoctorService(client)

	handler := routes.NewHandler(svc)
	router.POST("/getVideo", handler.GetVideos)
	router.Run()

}
