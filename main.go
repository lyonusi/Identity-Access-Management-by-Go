package main

import (
	"IAMbyGo/api"
	"IAMbyGo/middlewares"
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

	public := echo.New()

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
	middlewares := middlewares.NewMidWare(authService)

	// Echo instance
	public.Use(middleware.CORS())
	public.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Skipper: func(c echo.Context) bool {
			return true
		},
		SigningKey: []byte(service.TokenKey),
		AuthScheme: "Bearer",
		// TokenLookup: "header:" + scopeHeader,
	}))

	readonly := public.Group("/user")
	readonly.Use(middlewares.ValidateJWT, middlewares.CheckScopeFuncFactory("read"))

	// Routes - user
	readonly.GET("/hello", endpoint.Hello)
	readonly.GET("/getUserById", endpoint.GetUserById)
	readonly.GET("/listuser", endpoint.List)
	readonly.GET("/refreshtoken", endpoint.RefreshToken)

	u := public.Group("/test")
	u.Use(middlewares.ValidateJWT, middlewares.CheckScopeFuncFactory("test"))
	u.GET("/hello", endpoint.Hello)

	admin := public.Group("/admin")
	admin.Use(middlewares.ValidateJWT, middlewares.CheckScopeFuncFactory("read", "write"))
	// admin.Use(middleware.JWTWithConfig(middleware.JWTConfig{
	// 	SigningKey: []byte(service.TokenKey),
	// 	AuthScheme: "Bearer",
	// 	// TokenLookup: "header:" + scopeHeader,
	// }),)

	// Routes - admin
	admin.GET("/getUserById", endpoint.GetUserById)
	admin.POST("/createuser", endpoint.CreateUser)
	admin.GET("/listuser", endpoint.List)
	admin.POST("/update", endpoint.Update)
	admin.POST("/delete", endpoint.Delete)
	admin.GET("/refreshtoken", endpoint.RefreshToken)
	admin.POST("/setscope", endpoint.SetUserScope)
	admin.GET("/getscope", endpoint.GetUserScope)
	admin.GET("/listuserbyscope", endpoint.ListUserbyScope)
	admin.POST("/deleteuserscope", endpoint.DeleteUserScope)

	// Routes - public
	// public.GET("/hello", endpoint.Hello)
	public.POST("/login", endpoint.LogIn)
	public.POST("/emaillogin", endpoint.EmailLogIn)

	public.File("/", "./frontend/index.html")
	public.POST("/login-form", endpoint.LoginForm)

	// Start server
	public.Logger.Fatal(public.Start(":1323"))
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
