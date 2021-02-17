package auth

import (
	"context"
	"fmt"
	"gowithmongo/database"
	"gowithmongo/seeder"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

func setup() *mongo.Database {
	viper.SetConfigFile(`../../../config/config.json`)
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	dbURL := viper.GetString(`test.database.url`)
	dbName := viper.GetString(`test.database.name`)
	dbUser := viper.GetString(`test.database.username`)
	dbPass := viper.GetString(`test.database.password`)
	db := database.ConnectDB(dbURL, dbUser, dbPass, dbName)

	return db
}

func tearDown(db mongo.Database) {
	db.Drop(context.TODO())
}

// TestMain overrides main test
// We add data seeding here, so we don't have to execute data seeder for each test
func TestMain(m *testing.M) {
	// setup database
	db := setup()

	//seed user
	seeder.SeedUser(*db)

	code := m.Run()

	tearDown(*db)
	os.Exit(code)

}

//TestAuthLogin check response body and code for login with valid credential
func TestAuthLogin(t *testing.T) {
	db := setup()

	e := echo.New()
	email := "admin@admin.com"
	password := "admin"

	payload := fmt.Sprintf(`
		{
			"email": "%s",
			"password": "%s"
		}
	`, email, password)

	request := httptest.NewRequest("POST", "/auth/login", strings.NewReader(payload))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)

	authRepo := NewRepository(*db)
	h := &Handler{
		authRepo: authRepo,
	}

	if assert.NoError(t, h.authLogin(c)) {
		assert.Contains(t, rec.Body.String(), "access_token")
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

//TestInvalidLogin checks response code for login with invalid credential
func TestInvalidLogin(t *testing.T) {
	db := setup()

	e := echo.New()
	email := "admin@admin.com"
	password := "password"

	payload := fmt.Sprintf(`
		{
			"email": "%s",
			"password": "%s"
		}
	`, email, password)

	request := httptest.NewRequest("POST", "/auth/login", strings.NewReader(payload))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(request, rec)

	authRepo := NewRepository(*db)
	h := &Handler{
		authRepo: authRepo,
	}

	if assert.NoError(t, h.authLogin(c)) {
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	}
}
