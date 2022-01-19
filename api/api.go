package api

import (
	"IAMbyGo/service"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type Api interface {
	// Hello(c echo.Context) error
	// Hello(w http.ResponseWriter, req *http.Request)
	CreateUser(c echo.Context) error
	GetUserById(c echo.Context) error
	List(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	LogIn(c echo.Context) error
	EmailLogIn(c echo.Context) error
	RefreshToken(c echo.Context) error
	LoginForm(c echo.Context) error
	SetUserScope(c echo.Context) error
	DeleteUserScope(c echo.Context) error
	GetUserScope(c echo.Context) error
	ListUserbyScope(c echo.Context) error
}
type api struct {
	userService service.User
	authService service.Auth
}

func NewApi(userService service.User, authService service.Auth) Api {
	return &api{
		userService: userService,
		authService: authService,
	}
}

type loginResponse struct {
	UserID    string   `json:"userID"`
	UserName  string   `json:"userName"`
	UserEmail string   `json:"userEmail"`
	Token     string   `json:"token"`
	Scope     []string `json:"scope"`
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// Handler
// func (a *api) Hello(c echo.Context) error {
// 	token := c.Request().Header[echo.HeaderAuthorization][0]
// 	fmt.Println(token)
// 	// return c.String(http.StatusOK, "Hello, World!")
// 	// token, _ := a.authService.Sign("123123")
// 	result, err := a.authService.Validate(token)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal Error: %s", err.Error()))
// 	}
// 	return c.JSON(http.StatusOK, result)
// }
// func (a *api) Hello(w http.ResponseWriter, req *http.Request) {
// 	token := req.Header.Get("Authorization")
// 	fmt.Println(token)
// 	// return c.String(http.StatusOK, "Hello, World!")
// 	// token, _ := a.authService.Sign("123123")
// 	// result, err := a.authService.Validate(token)
// 	// if err != nil {
// 	// 	return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal Error: %s", err.Error()))
// 	// }
// }

func (a *api) CreateUser(c echo.Context) error {
	userName := c.FormValue("name")
	userEmail := c.FormValue("email")
	password := c.FormValue("password")

	// fmt.Println(userName)
	// fmt.Println(password)
	err := a.userService.CreateUser(userName, userEmail, password)
	if err != nil {
		// fmt.Println(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal Error: %s", err.Error()))
	}
	return c.String(http.StatusOK, fmt.Sprintf("New user ["+userName+"] created with email ["+userEmail+"]"))
}

func (a *api) GetUserById(c echo.Context) error {
	id := c.QueryParam("userID")
	user, err := a.userService.GetUserByID(id)
	if err != nil {
		// fmt.Println(err.Error())
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Internal Error: %s", err.Error()))
	}
	// fmt.Println(user)
	return c.JSON(http.StatusOK, user)
}

func (a *api) List(c echo.Context) error {
	userList, err := a.userService.List()
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Internal Error: %s", err.Error()))
	}
	return c.JSON(http.StatusOK, userList)
}

func (a *api) Update(c echo.Context) error {
	name := c.FormValue("name")
	password := c.FormValue("password")
	email := c.FormValue("email")
	id := c.FormValue("userID")

	switch field := c.QueryParam("field"); field {
	case "user":
		// fmt.Println("...Field accepted - ", field)
		err1 := a.userService.UpdateUser(service.UserInfo{
			UserID:    id,
			UserName:  name,
			UserEmail: email,
		})
		if err1 != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal Error: %s", err1.Error()))
		} else {
			return c.String(http.StatusOK, fmt.Sprintf("User ID "+id+" updated"))
		}
	case "password":
		// fmt.Println("...Field accepted - ", field)
		err2 := a.userService.UpdatePassword(id, password)
		if err2 != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal Error: %s", err2.Error()))
		} else {
			return c.String(http.StatusOK, "User ID "+id+" - password updated.")

		}
	default:
		return echo.NewHTTPError(http.StatusNotFound, "Internal Error: Field Not Found")
	}
}

func (a *api) Delete(c echo.Context) error {
	id := c.FormValue("id")
	name, err := a.userService.DeleteUser(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal Error: %s", err.Error()))
	}
	return c.String(http.StatusOK, fmt.Sprintf("User deleted: "+name+", user ID = "+id))
}

func (a *api) LogIn(c echo.Context) error {
	name := c.FormValue("username")
	password := c.FormValue("password")
	// fmt.Println("username: ", name)
	// fmt.Println("password: ", password)
	_, err := a.authService.LogIn(name, password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal Error: %s", err.Error()))
	}
	userID, _, err := a.userService.GetPasswordByName(name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal Error: %s", err.Error()))
	}
	token, err := a.authService.Sign(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal Error: %s", err.Error()))
	}
	userInfo, err := a.userService.GetUserByID(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal Error: %s", err.Error()))
	}
	usernameLoginResponse := loginResponse{
		UserID:    userInfo.UserID,
		UserName:  userInfo.UserName,
		UserEmail: userInfo.UserEmail,
		Token:     token,
		Scope:     userInfo.Scope,
	}

	return c.JSON(http.StatusOK, usernameLoginResponse)
}

func (a *api) EmailLogIn(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	_, err := a.authService.EmailLogIn(email, password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal Error: %s", err.Error()))
	}
	userID, _, err := a.userService.GetPasswordByEmail(email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal Error: %s", err.Error()))
	}

	// fmt.Println(userID)

	userInfo, err := a.userService.GetUserByID(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal Error: %s", err.Error()))
	}

	token, err := a.authService.Sign(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal Error: %s", err.Error()))
	}

	// fmt.Println(userInfo)

	emailLoginResponse := loginResponse{
		UserID:    userInfo.UserID,
		UserName:  userInfo.UserName,
		UserEmail: userInfo.UserEmail,
		Token:     token,
	}

	return c.JSON(http.StatusOK, emailLoginResponse)
}

func (a *api) RefreshToken(c echo.Context) error {
	token := strings.TrimPrefix(c.Request().Header[echo.HeaderAuthorization][0], "Bearer ")
	newToken, err := a.authService.RefreshToken(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal Error: %s", err.Error()))
	}
	return c.String(http.StatusOK, newToken)
}

func (a *api) LoginForm(c echo.Context) error {
	name := c.FormValue("username")
	password := c.FormValue("password")
	// fmt.Println("username: ", name)
	// fmt.Println("password: ", password)
	_, err := a.authService.LogIn(name, password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal Error: %s", err.Error()))
	}
	userID, _, err := a.userService.GetPasswordByName(name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal Error: %s", err.Error()))
	}
	token, err := a.authService.Sign(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal Error: %s", err.Error()))
	}
	t := &Template{
		templates: template.Must(template.ParseGlob("./frontend/admin.html")),
	}

	c.Echo().Renderer = t
	return c.Render(http.StatusOK, "hello", map[string]interface{}{
		"username": name,
		"userID":   userID,
		"token":    token,
	})
}

func (a *api) SetUserScope(c echo.Context) error {
	id := c.FormValue("userID")
	scope := c.FormValue("scope")

	err := a.userService.SetScopeByID(id, scope)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal Error: %s", err.Error()))
	}
	return c.String(http.StatusOK, fmt.Sprintf("User id ["+id+"] updated with ["+scope+"] scope"))
}
func (a *api) DeleteUserScope(c echo.Context) error {
	id := c.FormValue("userID")
	scope := c.FormValue("scope")

	err := a.userService.DeleteScopeByID(id, scope)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal Error: %s", err.Error()))
	}
	return c.String(http.StatusOK, fmt.Sprintf("Deleted ["+scope+"] scope for userID ["+id+"]"))
}
func (a *api) GetUserScope(c echo.Context) error {
	id := c.FormValue("userID")
	userScope, err := a.userService.ListScopeByID(id)
	if err != nil {
		// fmt.Println(err.Error())
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Internal Error: %s", err.Error()))
	}
	// fmt.Println(user)
	return c.JSON(http.StatusOK, userScope)
}
func (a *api) ListUserbyScope(c echo.Context) error {
	scope := c.FormValue("scope")

	userList, err := a.userService.ListUserByScope(scope)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Internal Error: %s", err.Error()))
	}
	return c.JSON(http.StatusOK, userList)
}
