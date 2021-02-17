package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// ConnectDB init function for database
func ConnectDB(dbURL string, dbUser string, dbPass string, dbName string) *mongo.Database {
	option := options.Client().ApplyURI(dbURL)
	if len(dbUser) != 0 {
		credential := options.Credential{
			AuthSource: dbName,
			Username:   dbUser,
			Password:   dbPass,
		}
		option = option.SetAuth(credential)
	}

	client, err := mongo.NewClient(option)
	if err != nil {
		panic(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	db := client.Database(dbName)

	return db
}
