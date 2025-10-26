package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nyxoy77/B2C_YouTube_Doctor/routes"
)

func main() {
	router := gin.Default()

	router.POST("/getVideo", routes.GetVideos)

	err := godotenv.Load()
	if err != nil {
		fmt.Printf("error loading the env : %v", err)
		return
	}
	router.Run()

}
