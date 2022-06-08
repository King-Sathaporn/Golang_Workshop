package api

import (
	"main/db"

	"github.com/gin-gonic/gin"
)

// Setup - call this function to setup the router
func Setup(router *gin.Engine) {
	db.SetupDB()
	SetupAuthenAPI(router)
	SetupProductAPI(router)
	SetupTransactionAPI(router)
}
