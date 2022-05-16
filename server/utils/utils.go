package utils

import (
	"encoding/json"
	"fmt"

	"github.com/brianvoe/sjwt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PrintJson(s gin.H) {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}

func ParseJWT(jwt string) sjwt.Claims {
	claims, err := sjwt.Parse(jwt)
	if err != nil {
		panic(err)
	}

	return claims
}

// Gets the user id from the current session as an ObjectID object
func GetUserId(session sessions.Session) primitive.ObjectID {
	objectID, err := primitive.ObjectIDFromHex(session.Get("userId").(string))
	if err != nil {
		panic(err)
	}

	return objectID
}
