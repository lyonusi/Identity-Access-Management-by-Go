package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// type UserInfo struct {
// 	UserID   string `json:"userID"`
// 	UserName string `json:"usersName"`
// }

var TokenKey = []byte("STeZg1g5IEwyGlD/5fiBjrJ+WtXDlU2SxKMWlJuwAAM=")

type Auth interface {
	LogIn(userName string, password string) (string, error)
	Sign(userID string) (string, error)
	RefreshToken(token string) (string, error)
	// Validate(token string) (userID string, err error)
}

type auth struct {
	userService User
}

func NewAuth(userService User) Auth {
	return &auth{
		userService: userService,
	}
}

func (a *auth) LogIn(userName string, password string) (string, error) {
	_, hashed, err := a.userService.GetUserPassword(userName)
	if err != nil {
		return "", fmt.Errorf("service.LogIn: %s", err.Error())
	}
	match := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if match == nil {
		return fmt.Sprintf("User " + userName + " logged in"), nil
	} else {
		return "", fmt.Errorf("log in failed")
	}
}

func (a *auth) Sign(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Minute * 30).Unix(),
	})
	tokenString, err := token.SignedString(TokenKey)
	if err != nil {
		return "", fmt.Errorf("service.Sign: %s", err.Error())
	}
	// fmt.Println(tokenString)
	return tokenString, nil
}

// func (a *auth) Validate(tokenString string) (userID string, err error) {
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		return TokenKey, nil
// 	})
// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		fmt.Println(claims["userID"])
// 		return claims["userID"].(string), nil
// 	} else {
// 		return "", fmt.Errorf("validate: %s", err.Error())
// 	}

// }

func (a *auth) RefreshToken(tokenString string) (string, error) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return "", nil
	})
	userID := token.Claims.(jwt.MapClaims)["userID"]
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Minute * 30).Unix(),
	})
	newTokenString, err := newToken.SignedString(TokenKey)
	if err != nil {
		return "", fmt.Errorf("service.refresh: %s", err.Error())
	}
	return newTokenString, nil
}
