package main

import (
	"flag"
	"fmt"
	"gowithmongo/database"
	mdl "gowithmongo/middleware"
	"gowithmongo/src/modules/auth"
	"gowithmongo/src/modules/user"
	"net/http"
	"os"

	_ "gowithmongo/docs"
	"gowithmongo/seeder"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	viper.SetConfigFile(`./config/config.json`)
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

}

func handleArgs(db *mongo.Database) {
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		switch args[0] {
		case "seed":
			// seed users
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
	// init database
	dbURL := viper.GetString(`database.url`)
	dbName := viper.GetString(`database.name`)
	dbUser := viper.GetString(`database.username`)
	dbPass := viper.GetString(`database.password`)
	db := database.ConnectDB(dbURL, dbUser, dbPass, dbName)

	handleArgs(db)

	appName := viper.GetString(`app.name`)
	appHost := viper.GetString(`app.host`)
	appPort := viper.GetString(`app.port`)

	e := echo.New()
	appMiddleware := mdl.InitAppMiddleware(appName)

	e.Use(appMiddleware.CORS)

	e.GET("/", hello)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	auth.NewAuthHandler(e, *db)
	user.NewUserHandler(e, *db)

	e.Logger.Fatal(e.Start(fmt.Sprintf(`%s:%s`, appHost, appPort)))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World")
}
