package main

import (
	"database/sql"
	"flag"
	"user-auth/internal/handlers"
	"user-auth/internal/models"
	"user-auth/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	port := flag.String("port", ":8000", "HTTP port")
	dsn := flag.String("dsn", "root:password@tcp(mysql_user:3306)/user", "MySQL data source name")
	key := flag.String("key", "secret", "Private key JWT")
	flag.Parse()
	e := echo.New()
	e.Validator = utils.NewValidator()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db, err := openDB(*dsn)
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer db.Close()

	h := handlers.UserHandler{UserModel: models.UserModel{DB: db}}
	h.SetKey(*key)
	authGroup := e.Group("/auth")
	authGroup.POST("/signup", h.Signup)
	authGroup.POST("/login", h.Login)

	userGroup := e.Group("/users")
	userGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(*key),
	}))
	userGroup.GET("/", h.GetAll)
	userGroup.PUT("/:id", h.Update)
	userGroup.DELETE("/:id", h.Delete)
	userGroup.GET("/:id", h.Get)
	e.Logger.Fatal(e.Start(*port))
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
