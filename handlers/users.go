package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"user-auth/models"
	"user-auth/utils"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
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
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err}
	}
	if u.Username == "" || u.Email == "" || u.Password == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid data"}
	}
	pw, err := utils.HashPassword(u.Password)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err}
	}
	u.Password = pw
	id, err := h.UserModel.Insert(*u)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err}
	}

	return c.JSON(http.StatusCreated, fmt.Sprintf("id:%d", id))
}

func (h *UserHandler) Login(c echo.Context) (err error) {
	u := new(models.User)
	if err = c.Bind(u); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	user, err := h.UserModel.GetByName(u.Username)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}

	if !utils.CheckPasswordHash(u.Password, user.Password) {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "wrong password"}
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	log.Printf("userID=%d", user.Id)
	claims["id"] = strconv.Itoa(user.Id)
	claims["is_admin"] = user.IsAdmin
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	user.Token, err = token.SignedString([]byte(Key))
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}
	user.Password = ""
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Get(c echo.Context) (err error) {
	currID := userIDFromToken(c)
	id := c.Param("id")
	log.Printf("CurrID=%s, IDParam=%s", currID, id)
	idConv, err := strconv.Atoi(id)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err}
	}
	user, err := h.UserModel.Get(idConv)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}
	if currID != id {
		user.Email = ""
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetAll(c echo.Context) (err error) {
	isAdmin := isAdminFromToken(c)
	if !isAdmin {
		return &echo.HTTPError{Code: http.StatusForbidden, Message: "forbidden"}
	}
	users, err := h.UserModel.List()
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}
	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) Update(c echo.Context) (err error) {
	currID := userIDFromToken(c)
	id := c.Param("id")
	if currID != id {
		return &echo.HTTPError{Code: http.StatusForbidden, Message: "forbidden"}
	}
	u := new(models.User)
	if err = c.Bind(u); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}
	if u.Id, err = strconv.Atoi(id); err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: err}
	}
	err = h.UserModel.Update(*u)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}
	return c.JSON(http.StatusOK, "updated")
}

func (h *UserHandler) Delete(c echo.Context) (err error) {
	currID := userIDFromToken(c)
	id := c.Param("id")
	if currID != id {
		return &echo.HTTPError{Code: http.StatusForbidden, Message: "forbidden"}
	}
	idConv, _ := strconv.Atoi(id)
	err = h.UserModel.Delete(idConv)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: err}
	}
	return c.JSON(http.StatusOK, "deleted")
}

func userIDFromToken(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["id"].(string)
}

func isAdminFromToken(c echo.Context) bool {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["is_admin"].(bool)
}
