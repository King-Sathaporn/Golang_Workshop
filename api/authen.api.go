package api

import (
	"main/db"
	"main/interceptor"
	"main/model"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SetupAuthenAPI(router *gin.Engine) {
	authenAPI := router.Group("/api/v2")
	{
		authenAPI.POST("/login", login)
		authenAPI.POST("/register", register)
	}
}

func login(c *gin.Context) {
	var user model.User
	if c.ShouldBind(&user) == nil {

		var queryUser model.User

		if err := db.GetDB().First(&queryUser, "username = ?", user.Username).Error; err != nil {

			c.JSON(401, gin.H{"result": "error", "error": err})

		} else if !checkPassword(user.Password, queryUser.Password) {
			//? !checkPassword(user.Password, queryUser.Password) equals to checkPassword(user.Password, queryUser.Password) == false

			c.JSON(401, gin.H{"result": "Invalid username or password"})

		} else {

			token := interceptor.JWTSign(queryUser)

			c.JSON(200, gin.H{
				"result": "success",
				"token":  token,
			})
		}
	} else {
		c.JSON(401, gin.H{"result": "unable to bind data."})
	}
}

func register(c *gin.Context) {

	var user model.User

	if c.ShouldBind(&user) == nil {
		user.Password, _ = hashPassword(user.Password)
		user.CreatedAt = time.Now()

		if err := db.GetDB().Create(&user).Error; err != nil {
			c.JSON(500, gin.H{
				"result": "error",
				"error":  err,
			})
		} else {
			c.JSON(200, gin.H{
				"result": "success",
				"token":  user,
			})
		}
	} else {
		c.JSON(401, gin.H{"status": "unable to bind data."})
	}

}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
