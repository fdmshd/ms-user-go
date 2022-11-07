package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"user-auth/internal/handlers"
	"user-auth/internal/models"
	"user-auth/internal/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const testKey = "test"

var CorrectUserJSON = `{
    "username":"Test",
    "email":"test@mail.ru",
    "password": "password"
}`

var RandomJSON = `{
	"field1":"Test",
    "field2":"test@mail.ru",
    "field3": "password"
}
`

func TestSignup(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(CorrectUserJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1))
	h := handlers.UserHandler{UserModel: models.UserModel{DB: db}, Validator: utils.NewValidator()}
	h.SetKey(testKey)

	if assert.NoError(t, h.Signup(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestSignupWrongJSON(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(RandomJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1))
	h := handlers.UserHandler{UserModel: models.UserModel{DB: db}, Validator: utils.NewValidator()}
	h.SetKey(testKey)
	err = h.Signup(c)
	expectedErr := echo.ErrBadRequest
	if assert.Error(t, err) {
		assert.IsType(t, expectedErr, err)
	}
}

func TestLogin(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(CorrectUserJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	columns := []string{"id", "username", "email", "is_admin", "password"}
	mock.ExpectQuery("SELECT(.+) FROM users WHERE username = (.+)").
		WillReturnRows(sqlmock.NewRows(columns).
			FromCSVString("1,Test,test@mail.ru,FALSE,$2a$14$5u2diVJLjdBWITCiXSO9SOj/YiPtL67BxyXP8lVdPbfEaGEM9b3vO"))
	h := handlers.UserHandler{UserModel: models.UserModel{DB: db}, Validator: utils.NewValidator()}
	h.SetKey("test")

	if assert.NoError(t, h.Login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestGetUser(t *testing.T) {
	e := echo.New()
	e.Validator = utils.NewValidator()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")
	user := models.User{Id: 1, Username: "Test", Email: "test@mail.ru"}

	c.Set("user", handlers.NewUserToken(&user))
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "username", "email"}).
		AddRow(user.Id, user.Username, user.Email)
	mock.ExpectQuery("SELECT(.+) FROM users WHERE id = (.+)").
		WillReturnRows(rows)
	h := handlers.UserHandler{UserModel: models.UserModel{DB: db}}
	h.SetKey(testKey)

	expectedJson, _ := json.Marshal(user)
	if assert.NoError(t, h.Get(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), string(expectedJson))
	}
}

func TestGetOtherUser(t *testing.T) {
	e := echo.New()
	e.Validator = utils.NewValidator()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("2")
	user1 := models.User{Id: 1, Username: "Test1", Email: "test1@mail.ru"}
	user2 := models.User{Id: 2, Username: "Test2", Email: "test2@mail.ru"}
	c.Set("user", handlers.NewUserToken(&user1))
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"id", "username", "email"}).
		AddRow(user2.Id, user2.Username, user2.Email)
	mock.ExpectQuery("SELECT(.+) FROM users WHERE id = (.+)").
		WillReturnRows(rows)
	h := handlers.UserHandler{UserModel: models.UserModel{DB: db}}
	h.SetKey(testKey)
	user2.Email = ""
	expectedJson, _ := json.Marshal(user2)
	if assert.NoError(t, h.Get(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), string(expectedJson))
	}
}
