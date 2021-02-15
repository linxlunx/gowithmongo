package auth

import (
	"context"
	"errors"
	"fmt"
	"gowithmongo/src/modules/user"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

// Login check user exist or not
func (r *Repository) Login(email string, password string) (user.User, error) {
	var userDetail user.User

	err := r.coll.FindOne(context.TODO(), bson.M{"email": email}).Decode(&userDetail)

	if err != nil {
		fmt.Println(err)
		return userDetail, errors.New("Wrong Email/Password")
	}

	errPasswd := bcrypt.CompareHashAndPassword([]byte(userDetail.Password), []byte(password))
	if errPasswd != nil {
		return userDetail, errors.New("Wrong Email/Password")
	}

	return userDetail, nil
}
