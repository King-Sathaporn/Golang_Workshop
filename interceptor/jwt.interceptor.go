package interceptor

import (
	"fmt"
	"main/model"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var secretKey = []byte("y#Y8g5$w&F3##Nry")

func JWTSign(payload model.User) string {
	atClaims := jwt.MapClaims{}
	atClaims["id"] = payload.ID
	atClaims["username"] = payload.Username
	atClaims["Level"] = payload.Level
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, _ := at.SignedString(secretKey)

	return token
}

func JWTVerify(c *gin.Context) {
	tokenString := strings.Split(c.Request.Header["Authorization"][0], " ")[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
		staffID := fmt.Sprintf("%v", claims["id"])
		username := fmt.Sprintf("%v", claims["jwt_username"])
		level := fmt.Sprintf("%v", claims["jwt_level"])
		c.Set("jwt_staff_id", staffID)
		c.Set("jwt_username", username)
		c.Set("jwt_level", level)
		c.Next()
	} else {
		c.JSON(200, gin.H{
			"result":  "error",
			"message": "Invalid Token",
			"error":   err,
		})
		c.Abort()
	}

}
