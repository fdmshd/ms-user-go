package handlers_test

import (
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
	e.Validator = utils.NewValidator()
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
	h := handlers.UserHandler{UserModel: models.UserModel{DB: db}}
	h.SetKey("test")

	if assert.NoError(t, h.Signup(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestSignupWrongJSON(t *testing.T) {
	e := echo.New()
	e.Validator = utils.NewValidator()
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
	h := handlers.UserHandler{UserModel: models.UserModel{DB: db}}
	h.SetKey("test")
	err = h.Signup(c)
	expectedErr := echo.ErrBadRequest
	if assert.Error(t, err) {
		assert.IsType(t, expectedErr, err)
	}
}

func TestLogin(t *testing.T) {
	e := echo.New()
	e.Validator = utils.NewValidator()
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
	h := handlers.UserHandler{UserModel: models.UserModel{DB: db}}
	h.SetKey("test")

	if assert.NoError(t, h.Login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
