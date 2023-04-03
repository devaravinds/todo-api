package routes

import (
	controller "dexlock.com/todo-project/controllers"
	middleware "dexlock.com/todo-project/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/user", controller.GetUser())
	incomingRoutes.PATCH("user/active-status", controller.SetAsActive())
}