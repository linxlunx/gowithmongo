package main

import (
	"context"
	"flag"
	"fmt"
	mdl "gowithmongo/middleware"
	"gowithmongo/src/modules/auth"
	"gowithmongo/src/modules/user"
	"net/http"
	"os"
	"time"

	_ "gowithmongo/docs"
	"gowithmongo/seeder"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func init() {
	viper.SetConfigFile(`./config/config.json`)
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

}

// ConnectDB init function for database
func ConnectDB() *mongo.Database {
	// init database
	dbURL := viper.GetString(`database.url`)
	dbName := viper.GetString(`database.name`)
	dbUser := viper.GetString(`database.username`)

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

	db := client.Database(dbName)

	return db
}

func handleArgs() {
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		switch args[0] {
		case "seed":
			// seed users
			db := ConnectDB()
			seeder.SeedUser(*db)
		default:
			fmt.Println(`Available command: seed (command for seeding customer data)`)
		}
		os.Exit(0)
	}
}

// @title Go Restful API with Echo
// @version 1.0
// @description My Boilerplate with Echo

// @contact.name Linggar Primahastoko
// @contact.url http://linggar.asia
// @contact.email x@linggar.asia

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @host localhost:5000
// @BasePath /
func main() {
	handleArgs()

	appName := viper.GetString(`app.name`)
	appHost := viper.GetString(`app.host`)
	appPort := viper.GetString(`app.port`)

	e := echo.New()
	appMiddleware := mdl.InitAppMiddleware(appName)
	e.Use(appMiddleware.CORS)

	e.GET("/", hello)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	db := ConnectDB()
	auth.NewAuthHandler(e, *db)
	user.NewUserHandler(e, *db)

	e.Logger.Fatal(e.Start(fmt.Sprintf(`%s:%s`, appHost, appPort)))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World")
}
