package auth

import (
	"gowithmongo/helper/wrapper"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

// Handler receive repository passed from new handler
type Handler struct {
	authRepo *Repository
}

// NewAuthHandler is init handler for auth module
func NewAuthHandler(e *echo.Echo, db mongo.Database) {
	// register repository here
	authRepo := NewRepository(db)
	ah := &Handler{
		authRepo: authRepo,
	}
	authGroup := e.Group("/auth")
	{
		authGroup.POST("/login", ah.authLogin)
	}

}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

// authLogin godoc
// @Summary Auth Login
// @Description Auth Login
// @Tags auth
// @ID auth-login
// @Param user body LoginParam true "Auth Login"
// @Accept  json
// @Produce  json
// @Success 200 {object} tokenResponse
// @Router /auth/login [post]
func (h *Handler) authLogin(c echo.Context) error {

	u := new(LoginParam)
	if err := c.Bind(u); err != nil {
		return wrapper.Error(http.StatusBadRequest, err.Error(), c)
	}

	user, err := h.authRepo.Login(u.Email, u.Password)
	if err != nil {
		return wrapper.Error(http.StatusUnauthorized, err.Error(), c)
	}

	// init token
	token := jwt.New(jwt.SigningMethodHS256)

	// add data to claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID.Hex()
	claims["email"] = user.Email
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// sign token
	t, err := token.SignedString([]byte(viper.GetString(`app.secret`)))
	if err != nil {
		return wrapper.Error(http.StatusBadRequest, "Bad Request", c)
	}

	result := tokenResponse{
		AccessToken: t,
	}

	return c.JSON(http.StatusOK, result)
}
