package handlers

import (
	"fmt"
	"net/http"
	"time"
	"user-auth/models"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	UserModel models.UserModel
}

const (
	// Key (Should come from somewhere else).
	Key = "secret"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (h *Handler) Signup(c echo.Context) (err error) {
	u := &models.User{}
	if err = c.Bind(u); err != nil {
		return
	}
	if u.Username == "" || u.Email == "" || u.Password == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "invalid data"}
	}
	pw, err := HashPassword(u.Password)
	log.Print(pw)
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

func (h *Handler) Login(c echo.Context) (err error) {
	// Bind
	u := new(models.User)
	if err = c.Bind(u); err != nil {
		return
	}

	// Find user
	user, err := h.UserModel.GetByName(u.Username)
	if err != nil {
		return
	}

	if !CheckPasswordHash(u.Password, user.Password) {
		return c.JSON(http.StatusBadRequest, "wrong password")
	}

	//-----
	// JWT
	//-----

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.Id
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response
	user.Token, err = token.SignedString([]byte(Key))
	if err != nil {
		return err
	}

	user.Password = "" // Don't send password
	return c.JSON(http.StatusOK, user)
}
