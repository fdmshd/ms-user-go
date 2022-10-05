package main

import (
	"database/sql"
	"os"
	"user-auth/handlers"
	"user-auth/models"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	dsn := os.Getenv("db")
	db, err := openDB(dsn)
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer db.Close()

	h := handlers.Handler{UserModel: models.UserModel{DB: db}}
	e.POST("/signup", h.Signup)
	e.POST("/login", h.Login)
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
