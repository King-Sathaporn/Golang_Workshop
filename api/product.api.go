package api

import (
	"fmt"
	"main/db"
	"main/interceptor"
	"main/model"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupProductAPI(router *gin.Engine) {
	productAPI := router.Group("/api/v2")
	{
		productAPI.GET("/product", interceptor.JWTVerify, getProduct)
		productAPI.GET("/product/:id", interceptor.JWTVerify, getProductByID)
		productAPI.POST("/product", interceptor.JWTVerify, createProduct)
		productAPI.PUT("/product", interceptor.JWTVerify, editProduct)
	}
}

/*
func getProducts(c *gin.Context) {
	var product []model.Product
	db.GetDB().Find(&product) //query all product
	c.JSON(200, gin.H{"result": product})
}
*/

func getProduct(c *gin.Context) {
	var product []model.Product
	keyword := c.Query("keyword")
	if keyword != "" {
		keyword = fmt.Sprintf("%%%s%%", keyword) //golang use %% = %
		db.GetDB().Where("name LIKE ?", keyword).Find(&product)
	} else {
		db.GetDB().Find(&product)
	}
	c.JSON(200, gin.H{"result": product})
}

func getProductByID(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)
	var product model.Product
	db.GetDB().Where("id = ?", id).First(&product)
	c.JSON(200, gin.H{"result": product})
}

func createProduct(c *gin.Context) {
	product := model.Product{}
	product.Name = c.PostForm("name")
	product.Stock, _ = strconv.ParseInt(c.PostForm("stock"), 10, 64) // 10 is base, 64 is bit size
	product.Price, _ = strconv.ParseFloat(c.PostForm("price"), 64)

	product.CreatedAt = time.Now()
	db.GetDB().Create(&product)

	image, _ := c.FormFile("image")
	saveImage(image, &product, c)

	c.JSON(200, gin.H{"result": product})
}

func saveImage(image *multipart.FileHeader, product *model.Product, c *gin.Context) {
	if image != nil {
		runningDir, _ := os.Getwd()
		product.Image = image.Filename
		extension := filepath.Ext(image.Filename)
		fileName := fmt.Sprintf("%d%s", product.ID, extension)
		filepath := fmt.Sprintf("%s/uploaded/images/%s", runningDir, fileName)

		if fileExists(filepath) {
			os.Remove(filepath)
		}
		c.SaveUploadedFile(image, filepath)
		db.GetDB().Model(&product).Update("image", fileName)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func editProduct(c *gin.Context) {
	var product model.Product
	id, _ := strconv.ParseInt(c.PostForm("id"), 10, 32) //id for reference of product
	product.ID = uint(id)
	product.Name = c.PostForm("name")
	product.Stock, _ = strconv.ParseInt(c.PostForm("stock"), 10, 64)
	product.Price, _ = strconv.ParseFloat(c.PostForm("price"), 64)

	db.GetDB().Save(&product)
	image, _ := c.FormFile("image")
	saveImage(image, &product, c)
	c.JSON(200, gin.H{"result": product})
}
