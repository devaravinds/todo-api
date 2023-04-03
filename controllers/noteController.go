package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"dexlock.com/todo-project/database"
	"dexlock.com/todo-project/helpers"
	"dexlock.com/todo-project/models"
	"github.com/gin-gonic/gin"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var noteCollection *mongo.Collection = database.OpenCollection(database.Client, "note")

var (
	upgrader      = websocket.Upgrader{} // Use default options
	activeCounter int32                // Counter for active connections
	counter int
)


func CreateNote() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		userId := c.GetString("uid") 
		var note models.Note
		if err := c.BindJSON(&note); err!= nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(note)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		note.ID = primitive.NewObjectID()
		note.User_id = userId

		resultInsertionNumber, insertErr := noteCollection.InsertOne(ctx, note)
		if insertErr != nil{
			msg:= fmt.Sprintf("Note was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return
		}
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func GetNotes() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		userId := c.GetString("uid") 

		filter := bson.M{"user_id": userId}

		var notes []models.Note

		cursor, err := noteCollection.Find(ctx, filter)
		if err!=nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
			return
		}

		for cursor.Next(ctx) {
			var note models.Note
			err := cursor.Decode(&note)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to decode note"})
				return
			}
			notes = append(notes, note)
		}

		c.JSON(http.StatusOK,notes)
	}
}

func GetNote() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		noteId,err := primitive.ObjectIDFromHex(c.Param("note_id"))
		if err!= nil{
			c.JSON(http.StatusInternalServerError, "Error converting id")
			return
		}

		var note models.Note
		err = noteCollection.FindOne(ctx, bson.M{"_id":noteId}).Decode(&note)
		if err!=nil{
			c.JSON(http.StatusInternalServerError, "Not found")
			return
		}

		userId := note.User_id
		err = helpers.MatchUserTypeToUid(c,userId)
		if err!=nil{
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		fmt.Print(userId)

		c.JSON(http.StatusOK, note)
	}
}

func UpdateNote() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		noteId,err := primitive.ObjectIDFromHex(c.Param("note_id")) 

		if err!=nil{
			c.JSON(http.StatusInternalServerError, "error converting id")
		}

		var note models.Note

		err = noteCollection.FindOne(ctx, bson.M{"_id":noteId}).Decode(&note)
		if err!=nil{
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		userId := note.User_id
		err = helpers.MatchUserTypeToUid(c,userId)
		if err!=nil{
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		err = c.BindJSON(&note)
		if err!=nil{
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		validationErr := validate.Struct(note)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		update := bson.M{"$set": bson.M{
			"content": note.Content,
			"title":note.Title,
		}}

		result,err := noteCollection.UpdateOne(ctx,bson.M{"_id":noteId},update)
		if err!=nil{
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK,result)
	}
}

func DeleteNote() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		noteId,err := primitive.ObjectIDFromHex(c.Param("note_id")) 
		if err!=nil{
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		var note models.Note

		err = noteCollection.FindOne(ctx, bson.M{"_id":noteId}).Decode(&note)
		if err!=nil{
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		userId := note.User_id
		err = helpers.MatchUserTypeToUid(c,userId)
		if err!=nil{
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		result,err := noteCollection.DeleteOne(ctx,bson.M{"_id":noteId})
		if err!=nil{
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, result)
	}
}




