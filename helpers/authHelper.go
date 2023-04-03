package helpers

import (
	"errors"


	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"dexlock.com/todo-project/database"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func MatchUserTypeToUid(c *gin.Context, userId string) (err error){
	uid := c.GetString("uid")
	err=nil

	if uid != userId{
		err = errors.New("Unauthorized")
		return err
	}
	return err
}


