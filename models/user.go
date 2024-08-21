package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username  string             `bson:"username" json:"username"`
	FirstName string             `bson:"firstName" json:"firstName"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"password"`
}
