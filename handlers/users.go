package handlers

import (
	"fmt"
	"net/http"
	"time"
	"user-auth/models"
	"user-auth/utils"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserModel models.UserModel
}

const (
	// Key (Should come from somewhere else).
	Key = "secret"
)

func (h *UserHandler) Signup(c echo.Context) (err error) {
	u := &models.User{}
	if err = c.Bind(u); err != nil {
		return
	}
	if u.Username == "" || u.Email == "" || u.Password == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid data"}
	}
	pw, err := utils.HashPassword(u.Password)
	if err != nil {
		return fmt.Errorf("hasing password:%v", err)
	}
	u.Password = pw
	id, err := h.UserModel.Insert(*u)
	if err != nil {
		return
	}

	return c.JSON(http.StatusCreated, fmt.Sprintf("id:%d", id))
}

func (h *UserHandler) Login(c echo.Context) (err error) {
	u := new(models.User)
	if err = c.Bind(u); err != nil {
		return
	}

	user, err := h.UserModel.GetByName(u.Username)
	if err != nil {
		return
	}

	if !utils.CheckPasswordHash(u.Password, user.Password) {
		return c.JSON(http.StatusBadRequest, "wrong password")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.Id
	claims["is_admin"] = user.IsAdmin
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	user.Token, err = token.SignedString([]byte(Key))
	if err != nil {
		return err
	}
	user.Password = ""
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetAllUsers(c echo.Context) (err error) {
	isAdmin := isAdminFromToken(c)
	if !isAdmin {
		return c.JSON(http.StatusForbidden, "forbidden")
	}
	users, err := h.UserModel.List()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, users)
}

func isAdminFromToken(c echo.Context) bool {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["is_admin"].(bool)
}
