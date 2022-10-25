package main

import (
	"database/sql"
	"user-auth/internal/handlers"
	"user-auth/internal/models"
	"user-auth/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	e := echo.New()
	e.Validator = utils.NewValidator()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	dsn := "root:password@tcp(mysql_user:3306)/user"
	db, err := openDB(dsn)
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer db.Close()

	h := handlers.UserHandler{UserModel: models.UserModel{DB: db}}
	e.POST("/signup", h.Signup)
	e.POST("/login", h.Login)

	userGroup := e.Group("/users")
	userGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(handlers.Key),
	}))
	userGroup.GET("", h.GetAll)
	userGroup.PUT("/:id", h.Update)
	userGroup.DELETE("/:id", h.Delete)
	userGroup.GET("/:id", h.Get)
	e.Logger.Fatal(e.Start(":8000"))
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
