package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nyxoy77/B2C_YouTube_Doctor/routes"
)

func main() {
	router := gin.Default()

	router.POST("/getVideo", routes.GetVideos)

	router.Run()
}
