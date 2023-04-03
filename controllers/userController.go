package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"dexlock.com/todo-project/database"
	"dexlock.com/todo-project/helpers"
	"dexlock.com/todo-project/models"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New() 

func HashPassword(password string) string{
	hashedPassword,err := bcrypt.GenerateFromPassword([]byte(password),14)
	if err!=nil {
		log.Panic(err)
	}
	return string(hashedPassword)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string){
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err!=nil{
		msg = fmt.Sprintf("email or password is incorrect")
		check = false
	}
	return check, msg
}

func SignUp() gin.HandlerFunc{
	return func(c *gin.Context) {
		//creating a new context with timeout = 100s
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		
		//declaring a user object and parsing the request body into it
		var user models.User
		if err := c.BindJSON(&user); err!= nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//validating the data
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		//checking if the email exists in db
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err!=nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking email"})
			return
		}
		if count>0{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email already exists"})
			return
		}

		//hashing the password
		password := HashPassword(user.Password)
		user.Password = password

		//creating userIds
		user.ID = primitive.NewObjectID()

		user.Is_active = false

		//inserting the document into mongodb
		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil{
			msg:= fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return
		}
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func Login() gin.HandlerFunc{
	return func(c *gin.Context){
		//creating a new context with timeout = 100s
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		//declaring 2 objects user and foundUser
		var authRequest models.AuthRequest
		var foundUser models.User



		//parsing the request body into user object
		if err:= c.BindJSON(&authRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}

		validationErr := validate.Struct(authRequest)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		//trying to match the provided credentials with the database 
		err := userCollection.FindOne(ctx, bson.M{"email":authRequest.Email}).Decode(&foundUser)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error":"email or password is incorrect"})
			return
		}

		//verifying password
		passwordIsValid,msg := VerifyPassword(authRequest.Password, foundUser.Password)
		if !passwordIsValid{
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		
		//error handling if user was not found
		if foundUser.Email == ""{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}

		//generating and updating tokens
		token, _ := helpers.GenerateTokens(foundUser.Email, foundUser.First_name, foundUser.Last_name, foundUser.ID.Hex())



		var response models.Response
		response.Token = token

		c.JSON(http.StatusOK, response)
	}
}



func GetUser() gin.HandlerFunc{
	return func(c *gin.Context){
		userId,err := primitive.ObjectIDFromHex(c.GetString("uid"))
		if err!= nil{
			c.JSON(http.StatusInternalServerError, "Error getting Id")
		} 

		fmt.Print(userId)
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		
		var user models.User
		err = userCollection.FindOne(ctx, bson.M{"_id":userId}).Decode(&user)
		defer cancel()
		if err!=nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
			return
		}
		c.JSON(http.StatusOK,user)
	}
}

func SetAsActive() gin.HandlerFunc{
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId,err := primitive.ObjectIDFromHex(c.GetString("uid"))
		if err!= nil{
			c.JSON(http.StatusInternalServerError, "error getting id")
			return
		} 

		var user models.User
		userCollection.FindOne(ctx, bson.M{"_id":userId}).Decode(&user)
	
		update := bson.M{"$set": bson.M{
			"is_active": !user.Is_active,
		}}
		userCollection.UpdateOne(ctx, bson.M{"_id":userId},update)

		c.JSON(http.StatusOK,"Active status is:"+ strconv.FormatBool(!user.Is_active))
	}
}