package user

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Repository defines repository for auth module
// return collection
type Repository struct {
	coll *mongo.Collection
}

// NewRepository init function for repository method
func NewRepository(db mongo.Database) *Repository {
	coll := db.Collection("users")
	return &Repository{
		coll: coll,
	}
}

// CreateUser is function to add users
func (r *Repository) CreateUser(param *NewUser) (User, error) {
	var userDetail User

	password, _ := bcrypt.GenerateFromPassword([]byte(param.Password), 14)
	addUser := User{
		ID:       primitive.NewObjectID(),
		Email:    param.Email,
		Fullname: param.Fullname,
		Password: string(password),
		Role:     "user",
	}
	resp, err := r.coll.InsertOne(context.TODO(), addUser)
	if err != nil {
		return userDetail, errors.New("Error Creating User")
	}

	err = r.coll.FindOne(context.TODO(), bson.M{"_id": resp.InsertedID}).Decode(&userDetail)
	if err != nil {
		return userDetail, errors.New("Error Creating User")
	}

	return userDetail, nil
}

// FindAll is function to list all users
func (r *Repository) FindAll() (Users, error) {
	var users Users

	cursor, err := r.coll.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Println(err)
	}

	for cursor.Next(context.TODO()) {
		var userDetail User
		err := cursor.Decode(&userDetail)
		if err != nil {
			log.Println(err)
		}
		users = append(users, userDetail)
	}

	return users, nil
}

// FindByID is function to get user by id
func (r *Repository) FindByID(userID string) (User, error) {
	var userDetail User

	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return userDetail, errors.New("Cannot find any user")
	}

	err = r.coll.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&userDetail)
	if err != nil {
		return userDetail, errors.New("Cannot find any user")
	}

	return userDetail, nil
}
