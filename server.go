package main

import (
	"fmt"
	"main/api"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	var port = os.Getenv("PORT")
	router := gin.Default()
	router.Static("/images", "./uploaded/images")

	api.Setup(router)
	if port == "" {
		fmt.Println("No Port In Heroku")
		router.Run()
	} else {
		fmt.Println("Environment Port: " + port)
		router.Run(":" + port)
	}
}
