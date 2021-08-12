package main

import (
	"database/sql"
	"fmt"
	"helloworld/repo"
	"helloworld/service"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

var userService service.User

func main() {
	fmt.Println("Server started....")

	db, err := dbInit()
	if err != nil {
		log.Fatal(err)
	}
	userRepo := repo.NewUser(db)
	userService = service.NewUser(userRepo)

	// Echo instance
	e := echo.New()

	// Routes
	e.GET("/", hello)
	e.GET("/getUserById", getUserById)
	e.POST("/createuser", createUser)
	e.GET("/listuser", list)
	e.POST("/update", update)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func dbInit() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "./user.db")
	if err != nil {
		return nil, fmt.Errorf("main.dbInit %s", err.Error())
	}

	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS users (userId VARCHAR(36) PRIMARY KEY, name VARCHAR(50), password VARCHAR(50))")
	if err != nil {
		return nil, err
	}

	statement.Exec()
	return database, nil
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func createUser(c echo.Context) error {
	userName := c.FormValue("name")
	password := c.FormValue("password")
	// fmt.Println(userName)
	// fmt.Println(password)
	err := userService.CreateUser(userName, password)
	if err != nil {
		// fmt.Println(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal Error: %s", err.Error()))
	}
	return c.String(http.StatusOK, "Hello, World!")
}

func getUserById(c echo.Context) error {
	id := c.QueryParam("userID")
	user, err := userService.GetUserByID(id)
	if err != nil {
		// fmt.Println(err.Error())
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Internal Error: %s", err.Error()))
	}
	// fmt.Println(user)
	return c.JSON(http.StatusOK, user)
}

func list(c echo.Context) error {
	userList, err := userService.List()
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Internal Error: %s", err.Error()))
	}
	return c.JSON(http.StatusOK, userList)
}

func update(c echo.Context) error {
	name := c.FormValue("name")
	password := c.FormValue("password")
	// fmt.Println("...name & password received - ", name, password)
	id := c.QueryParam("userID")
	// fmt.Println("...ID accepted - ", id)
	switch field := c.QueryParam("field"); field {
	case "name":
		// fmt.Println("...Field accepted - ", field)
		err1 := userService.UpdateName(id, name)
		if err1 != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal Error: %s", err1.Error()))
		}
		return err1
	case "password":
		// fmt.Println("...Field accepted - ", field)
		err2 := userService.UpdatePassword(id, password)
		if err2 != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal Error: %s", err2.Error()))
		}
		return err2
	default:
		return echo.NewHTTPError(http.StatusNotFound, "Internal Error: Field Not Found")
	}
}
