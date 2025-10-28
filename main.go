package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nyxoy77/B2C_YouTube_Doctor/client"
	"github.com/nyxoy77/B2C_YouTube_Doctor/routes"
	"github.com/nyxoy77/B2C_YouTube_Doctor/service"
	"github.com/nyxoy77/B2C_YouTube_Doctor/utils"
)

func main() {
	router := gin.Default()
	logger := utils.GlobalLoggerInstance()
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
	client := client.NewDoctorClient(ytService, genAiClient, logger)
	svc := service.NewDoctorService(client, logger)

	handler := routes.NewHandler(svc)
	router.POST("/getVideo", handler.GetVideos)
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	serverChannel := make(chan error, 1)
	signalChannel := make(chan os.Signal, 1)

	go func() {
		serverChannel <- server.ListenAndServe()
	}()

	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	select {
	case err := <-serverChannel:
		log.Fatalf("error occured starting the server: %v", err)
	case signal := <-signalChannel:
		log.Printf("signal recieved: %v", signal)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("erro occured while gracefully shutting down %v", err)
			server.Close()
		}

		log.Println("Server shutdown gracefully")
	}
}
