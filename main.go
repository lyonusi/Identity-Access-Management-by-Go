package main

import (
	"IAMbyGo/api"
	"IAMbyGo/repo"
	"IAMbyGo/service"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var userService service.User
var authService service.Auth

func main() {
	fmt.Println("Server starting....")
	db, err := dbInit()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	var redisClient *redis.Client
	redisClient, err = repo.RedisClient()

	if err != nil {
		log.Fatal(err)
	}

	userRepo := repo.NewUser(db)
	tokenRepo := repo.NewUserToken(db)
	userScope := repo.NewUserScope(db)
	userService = service.NewUser(userRepo, redisClient, userScope)
	authService = service.NewAuth(userService, tokenRepo)
	endpoint := api.NewApi(userService, authService)

	// Echo instance
	e.Use(middleware.CORS())
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Skipper: func(c echo.Context) bool {
			return true
		},
		SigningKey: []byte(service.TokenKey),
		AuthScheme: "Bearer",
	}))

	g := e.Group("/admin")
	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(service.TokenKey),
		AuthScheme: "Bearer",
	}))

	// Routes - admin
	g.GET("/getUserById", endpoint.GetUserById)
	g.POST("/createuser", endpoint.CreateUser)
	g.GET("/listuser", endpoint.List)
	g.POST("/update", endpoint.Update)
	g.POST("/delete", endpoint.Delete)
	g.GET("/refreshtoken", endpoint.RefreshToken)
	g.POST("/setscope", endpoint.SetUserScope)
	g.GET("/getscope", endpoint.GetUserScope)
	g.GET("/listuserbyscope", endpoint.ListUserbyScope)
	g.POST("/deleteuserscope", endpoint.DeleteUserScope)

	// Routes - public
	// e.GET("/hello", endpoint.Hello)
	e.POST("/login", endpoint.LogIn)
	e.POST("/emaillogin", endpoint.EmailLogIn)

	e.File("/", "./frontend/index.html")
	e.POST("/login-form", endpoint.LoginForm)

	// Routes - public, custom middleward
	// http.HandleFunc("/hello,", middlewares.Use(endpoint.Hello, middlewares.ValidateJWT))

	// Start server
	e.Logger.Fatal(e.Start(":1323"))

}

func dbInit() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "./user.db")
	if err != nil {
		return nil, fmt.Errorf("main.dbInit %s", err.Error())
	}

	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS users (userId VARCHAR(36) PRIMARY KEY, name VARCHAR(50), email VARCHAR(50), password VARCHAR(50))")
	if err != nil {
		return nil, err
	}
	statement.Exec()
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS token (userId VARCHAR(36) PRIMARY KEY, token VARCHAR(50))")
	if err != nil {
		return nil, err
	}
	statement.Exec()
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS user_scope (userId VARCHAR(36), scope VARCHAR(50))")
	if err != nil {
		return nil, err
	}
	statement.Exec()
	statement, err = database.Prepare("CREATE INDEX index_userId_scope ON user_scope (userId, scope)")
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return nil, err
		}
	} else {
		statement.Exec()
	}
	return database, nil
}
