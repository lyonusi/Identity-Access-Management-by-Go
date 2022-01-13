package service

import (
	"IAMbyGo/repo"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var TokenKey = []byte("STeZg1g5IEwyGlD/5fiBjrJ+WtXDlU2SxKMWlJuwAAM=")

type Auth interface {
	LogIn(userName string, password string) (string, error)
	EmailLogIn(userEmail string, password string) (string, error)
	Sign(userID string) (string, error)
	RefreshToken(tokenString string) (string, error)
	// Validate(token string) (userID string, err error)
}

type auth struct {
	userService User
	tokenRepo   repo.Token
}

func NewAuth(userService User, tokenRepo repo.Token) Auth {
	return &auth{
		userService: userService,
		tokenRepo:   tokenRepo,
	}
}

func (a *auth) LogIn(userName string, password string) (string, error) {
	_, hashed, err := a.userService.GetPasswordByName(userName)
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

func (a *auth) EmailLogIn(userEmail string, password string) (string, error) {
	_, hashed, err := a.userService.GetPasswordByEmail(userEmail)
	if err != nil {
		return "", fmt.Errorf("service.EmailLogIn: %s", err.Error())
	}
	tempUser, err := a.userService.GetUserByEmail(userEmail)
	if err != nil {
		return "", fmt.Errorf("service.EmailLogIn: %s", err.Error())
	}
	match := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if match == nil {
		return fmt.Sprintf("User " + tempUser.UserName + " logged in with " + tempUser.UserEmail), nil
	} else {
		return "", fmt.Errorf("log in failed")
	}
}

func (a *auth) Sign(userID string) (string, error) {
	scope, err := a.userService.ListScopeByID(userID)
	if err != nil {
		return "", fmt.Errorf("service.Sign.GetScope: %s", err.Error())
	}
	// fmt.Println("service.Sign.GetScope:", scope.UserScope)
	tokenShort := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"scope":  scope.UserScope,
		"exp":    time.Now().Add(time.Minute * 30).Unix(),
	})
	tokenStringShort, err := tokenShort.SignedString(TokenKey)
	if err != nil {
		return "", fmt.Errorf("service.Sign: %s", err.Error())
	}
	// tokenLong := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"userID": userID,
	// 	"exp":    time.Now().Add(time.Hour * 24 * 30).Unix(),
	// })
	// tokenStringLong, err := tokenLong.SignedString(TokenKey)
	// if err != nil {
	// 	return "", fmt.Errorf("service.Sign: %s", err.Error())
	// }
	// err = a.tokenRepo.CreateToken(userID, tokenStringLong)
	// if err != nil {
	// 	return "", fmt.Errorf("service.Sign: %s", err.Error())
	// }
	// fmt.Println("~~~~~short and long generated and long saved!~~~~~", tokenStringShort, tokenStringLong)
	return tokenStringShort, nil
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
	userID := fmt.Sprintf("%v", token.Claims.(jwt.MapClaims)["userID"])
	scope, err := a.userService.ListScopeByID(userID)
	if err != nil {
		return "", fmt.Errorf("service.Sign.GetScope: %s", err.Error())
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"scope":  scope.UserScope,
		"exp":    time.Now().Add(time.Minute * 30).Unix(),
	})
	newTokenString, err := newToken.SignedString(TokenKey)
	if err != nil {
		return "", fmt.Errorf("service.refresh: %s", err.Error())
	}
	return newTokenString, nil
}
