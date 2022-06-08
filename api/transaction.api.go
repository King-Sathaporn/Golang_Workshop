package api

import (
	"main/db"
	"main/interceptor"
	"main/model"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupTransactionAPI(router *gin.Engine) {
	TransactionAPI := router.Group("/api/v2")
	{
		TransactionAPI.GET("/transaction", getTransactions)
		TransactionAPI.POST("/transaction", interceptor.JWTVerify, createTransaction)
	}
}

type TransactionResult struct {
	ID            uint
	Total         float64
	Paid          float64
	Change        float64
	PaymentType   string
	PaymentDetail string
	OrderList     string
	Staff         string
	CreateAt      time.Time
}

func getTransactions(c *gin.Context) {
	var result []TransactionResult
	db.GetDB().Debug().Raw("SELECT transactions.id, total, paid, change, payment_type, payment_detail, order_list, users.username as Staff, transactions.create_at FROM transactions join users on transactions.staff_id = users.id", nil).Scan(&result)
	c.JSON(200, result)
}

func createTransaction(c *gin.Context) {
	var transaction model.Transaction
	if err := c.ShouldBindJSON(&transaction); err == nil {
		transaction.StaffID = c.GetString("jwt_staff_id")
		transaction.CreateAt = time.Now()
		db.GetDB().Create(&transaction)
		c.JSON(200, gin.H{"message": "success", "data": transaction})
	} else {
		c.JSON(404, gin.H{"message": "create transaction failed"})
	}
}
