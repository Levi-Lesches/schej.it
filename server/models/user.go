package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty" binding:"required"`
	Email     string             `json:"email" bson:"email" binding:"required"`
	FirstName string             `json:"firstName" bson:"firstName" binding:"required"`
	LastName  string             `json:"lastName" bson:"lastName" binding:"required"`
	Picture   string             `json:"picture" bson:"picture" binding:"required"`
}
