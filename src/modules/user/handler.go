package user

import (
	"gowithmongo/helper/wrapper"
	mdl "gowithmongo/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

// Handler receive repository passed from new handler
type Handler struct {
	userRepo *Repository
}

// NewUserHandler is init handler for users module
func NewUserHandler(e *echo.Echo, db mongo.Database) {
	// register repository here
	userRepo := NewRepository(db)
	uh := &Handler{
		userRepo: userRepo,
	}

	userGroup := e.Group("/users")
	userGroup.Use(mdl.IsLoggedIn(db))
	{
		userGroup.GET("", uh.userList)
		userGroup.GET("/:id", uh.userDetail)
		userGroup.POST("/add", uh.userAdd, mdl.IsAdmin(db))

	}

}

// userAdd godoc
// @Summary User Add
// @Description User Add
// @Tags users
// @ID users-add
// @Accept  json
// @Produce  json
// @Param user body NewUser true "User Add"
// @Success 200 {object} User
// @Router /users/add [post]
// @Security ApiKeyAuth
func (h *Handler) userAdd(c echo.Context) error {
	u := new(NewUser)
	if err := c.Bind(u); err != nil {
		return wrapper.Error(http.StatusBadRequest, err.Error(), c)
	}

	userAdd, err := h.userRepo.CreateUser(u)
	if err != nil {
		return wrapper.Error(http.StatusBadRequest, err.Error(), c)
	}
	return c.JSON(http.StatusOK, userAdd)
}

// userList godoc
// @Summary User List
// @Description User List
// @Tags users
// @ID users-list
// @Accept  json
// @Produce  json
// @Success 200 {object} Users
// @Router /users [get]
// @Security ApiKeyAuth
func (h *Handler) userList(c echo.Context) error {
	users, _ := h.userRepo.FindAll()
	return c.JSON(http.StatusOK, users)
}

// userDetail godoc
// @Summary User Get By ID
// @Description User Get By ID
// @Tags users
// @ID users-detail
// @Accept  json
// @Produce  json
// @Param userID path string false "User ID"
// @Success 200 {object} User
// @Router /users/{userID} [get]
// @Security ApiKeyAuth
func (h *Handler) userDetail(c echo.Context) error {
	user, err := h.userRepo.FindByID(c.Param("id"))
	if err != nil {
		return wrapper.Error(http.StatusUnauthorized, err.Error(), c)
	}

	return c.JSON(http.StatusOK, user)
}
