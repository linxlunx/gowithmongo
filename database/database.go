package database

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// ConnectDB maintains database connection
func ConnectDB() *mongo.Client {
	// init database
	dbURL := viper.GetString(`database.url`)
	dbName := viper.GetString(`database.name`)
	dbUser := viper.GetString(`database.username`)
	fmt.Println(dbURL)

	option := options.Client().ApplyURI(dbURL)
	if len(dbUser) != 0 {
		dbPass := viper.GetString(`database.password`)
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

	return client
}
