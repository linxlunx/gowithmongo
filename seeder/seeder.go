package seeder

import (
	"context"
	"fmt"
	"gowithmongo/src/modules/user"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type newUser struct {
	Email    string
	Fullname string
	Password string
	Role     string
}

// SeedUser seeds initial users
func SeedUser(db mongo.Database) {
	// define collection
	userCount, err := db.Collection("users").CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		panic(err)
	}

	if userCount > 0 {
		fmt.Println("Please empty the users collection before running this script")
		os.Exit(0)
	}

	adminPassword, _ := bcrypt.GenerateFromPassword([]byte("admin"), 14)
	userPassword, _ := bcrypt.GenerateFromPassword([]byte("user"), 14)

	users := []user.User{
		{
			ID:       primitive.NewObjectID(),
			Email:    "admin@admin.com",
			Fullname: "Administrator",
			Password: string(adminPassword),
			Role:     "administrator",
		},
		{
			ID:       primitive.NewObjectID(),
			Email:    "user@user.com",
			Fullname: "User",
			Password: string(userPassword),
			Role:     "user",
		},
	}

	for _, v := range users {
		_, err = db.Collection("users").InsertOne(context.TODO(), v, nil)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("Successfully seed initial users")
}
