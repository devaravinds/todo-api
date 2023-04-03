package helpers

import (
	"fmt"
	"os"
	"time"
	jwt "github.com/dgrijalva/jwt-go"
)

type SignedDetails struct{
	Email 			string
	First_name 		string
	Last_name 		string
	Uid 			string
	jwt.StandardClaims
}

// var userCollection *mongo.Collection = database.OpenCollection(database.Client,"user")
var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateTokens(email string, firstName string, lastName string, uid string ) (signedToken string, err error){
	claims := &SignedDetails{
		Email: email,
		First_name: firstName,
		Last_name: lastName,
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour*time.Duration(24)).Unix(),

		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	return token, err
}



func ValidateToken(signedToken string) (claims *SignedDetails, msg string){
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token)(interface{}, error){
			return[]byte(SECRET_KEY), nil
		},
	)

	if err!=nil{
		msg = err.Error()
		return
	}

	claims,ok := token.Claims.(*SignedDetails)
	if !ok{
		msg = fmt.Sprintf("The token is invalid")
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix(){
		msg = fmt.Sprintf("Token is expired")
		msg = err.Error()
		return
	}
	return claims, msg
}