package service

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// type UserInfo struct {
// 	UserID   string `json:"userID"`
// 	UserName string `json:"usersName"`
// }

type Auth interface {
	LogIn(userName string, password string) (string, error)
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
