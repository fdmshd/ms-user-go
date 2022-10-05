package handlers

import (
	"fmt"
	"net/http"
	"time"
	"user-auth/models"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	UserModel models.UserModel
}

const (
	// Key (Should come from somewhere else).
	Key = "secret"
)

func (h *Handler) Signup(c echo.Context) (err error) {
	u := &models.User{}
	if err = c.Bind(u); err != nil {
		return
	}
	if u.Username == "" || u.Email == "" || u.Password == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid data"}
	}
	id, err := h.UserModel.Insert(*u)
	if err != nil {
		return
	}

	return c.JSON(http.StatusCreated, fmt.Sprintf("new user id = %d", id))
}

func (h *Handler) Login(c echo.Context) (err error) {
	// Bind
	u := new(models.User)
	if err = c.Bind(u); err != nil {
		return
	}

	// Find user

	//-----
	// JWT
	//-----

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = u.Id
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response
	u.Token, err = token.SignedString([]byte(Key))
	if err != nil {
		return err
	}

	u.Password = "" // Don't send password
	return c.JSON(http.StatusOK, u)
}
