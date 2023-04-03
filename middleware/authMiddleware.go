package middleware

import (
	"net/http"

	"dexlock.com/todo-project/helpers"
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc{
	return func(c *gin.Context){
		clientToken := c.Request.Header.Get("token")
		if clientToken == ""{
			c.JSON(http.StatusUnauthorized, gin.H{"error":"No authorization header provided"})
			c.Abort()
			return
		}
		claims, err := helpers.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"Error":err})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("uid", claims.Uid)

	}
}
