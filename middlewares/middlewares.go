package middlewares

import (
	"IAMbyGo/service"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
)

const SCOPE_HEADER string = "X-scope"

var r = regexp.MustCompile(`\[(.*)\]`)

type JWTInfo struct {
	UserID string
	Scope  []string
}
type midware struct {
	authService service.Auth
}

func NewMidWare(authService service.Auth) Midware {
	return &midware{
		authService: authService,
	}
}

type Midware interface {
	ValidateJWT(next echo.HandlerFunc) echo.HandlerFunc
	CheckScopeFuncFactory(scope ...string) func(next echo.HandlerFunc) echo.HandlerFunc
}

func (m *midware) ValidateJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		token = strings.ReplaceAll(token, "Bearer ", "")
		// fmt.Println("middleware get token = ", token)
		result, err := m.authService.Validate(token)
		if err != nil {
			return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("Internal Error: %s", err.Error()))
		}
		// fmt.Println(fmt.Sprintf("%v", result.Scope))
		scopeString := r.FindStringSubmatch(fmt.Sprintf("%v", result.Scope))[1]
		// fmt.Println(scopeString)
		c.Request().Header.Set(SCOPE_HEADER, fmt.Sprintf("%v", scopeString))
		return next(c)
	}
}

func (m *midware) CheckScopeFuncFactory(scope ...string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			scopeString := c.Request().Header.Get(SCOPE_HEADER)
			scopeArray := strings.Split(scopeString, " ")
			scopeMap := make(map[string]bool, len(scopeArray))
			for _, s := range scopeArray {
				scopeMap[s] = true
			}
			for _, s := range scope {
				if _, ok := scopeMap[s]; !ok {
					return echo.NewHTTPError(http.StatusForbidden, "forbidden")
				}
			}
			return next(c)
		}
	}
}
