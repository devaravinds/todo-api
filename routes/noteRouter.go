package routes

import(
	"dexlock.com/todo-project/controllers"
	"github.com/gin-gonic/gin"
	"dexlock.com/todo-project/middleware"
)

func NoteRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.POST("/notes", controllers.CreateNote())
	incomingRoutes.GET("/notes", controllers.GetNotes())
	incomingRoutes.GET("/note/:note_id", controllers.GetNote())
	incomingRoutes.PATCH("/note/:note_id", controllers.UpdateNote())
	incomingRoutes.DELETE("note/:note_id", controllers.DeleteNote())
}