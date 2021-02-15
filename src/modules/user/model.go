package user

import "go.mongodb.org/mongo-driver/bson/primitive"

// User model is base structure for user collection
type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"-" bson:"password"`
	Fullname string             `json:"fullname" bson:"fullname"`
	Role     string             `json:"role" bson:"role"`
}

// Users is structure list for user
type Users []User
