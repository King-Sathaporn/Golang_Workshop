package main

import (
	"main/api"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Static("/images", "./uploaded/images")

	api.Setup(router)
	router.Run(getPort())
}

func getPort() string{
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}
	return ":" + port
}