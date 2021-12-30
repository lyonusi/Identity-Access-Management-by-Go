package service

import (
	"IAMbyGo/repo"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserInfo struct {
	UserID    string `json:"userID"`
	UserName  string `json:"userName"`
	UserEmail string `json:"userEmail"`
}

type User interface {
	CreateUser(userName string, userEmail string, password string) error
	GetUserByID(userID string) (*UserInfo, error)
	GetUserByEmail(userEmail string) (*UserInfo, error)
	List() ([]*UserInfo, error)
	UpdateName(userID string, userName string) error
	UpdatePassword(userID string, password string) error
	DeleteUser(userID string) (string, error)
	GetPasswordByName(userName string) (userID string, password string, err error)
	GetPasswordByEmail(userEmail string) (userID string, password string, err error)
}

type tools interface {
	hashPassword(password string) (string, error)
}

type tool struct {
}

type user struct {
	privateMethods tools
	userRepo       repo.User
}

func NewUser(userRepo repo.User) User {
	return &user{
		userRepo:       userRepo,
		privateMethods: &tool{},
	}
}

func (u *user) CreateUser(userName string, userEmail string, password string) error {
	id := uuid.New().String()
	hashPassword, err := u.privateMethods.hashPassword(password)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("service.CreateUser.HashPassword: %s", err.Error())
	}
	err = u.userRepo.CreateUser(id, userName, userEmail, hashPassword)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("service.CreateUser: %s", err.Error())
	}
	return nil
}

func (u *user) GetUserByID(userID string) (*UserInfo, error) {
	returnedUser, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("service.GetUserByID: %s", err.Error())
	}
	tempUser := &UserInfo{
		UserID:    returnedUser.UserID,
		UserName:  returnedUser.UserName,
		UserEmail: returnedUser.UserEmail,
	}

	return tempUser, nil
}

func (u *user) GetUserByEmail(userEmail string) (*UserInfo, error) {
	returnedUser, err := u.userRepo.GetUserByEmail(userEmail)
	if err != nil {
		return nil, fmt.Errorf("service.GetUserByEmail: %s", err.Error())
	}
	tempUser := &UserInfo{
		UserID:    returnedUser.UserID,
		UserName:  returnedUser.UserName,
		UserEmail: returnedUser.UserEmail,
	}
	return tempUser, nil
}

func (u *user) GetPasswordByName(userName string) (userID string, password string, err error) {
	returnedUser, err := u.userRepo.GetUserByName(userName)
	if err != nil {
		return "", "", fmt.Errorf("service.GetPasswordByName: %s", err.Error())
	}
	return returnedUser.UserID, returnedUser.Password, nil
}

func (u *user) GetPasswordByEmail(userEmail string) (userID string, password string, err error) {
	returnedUser, err := u.userRepo.GetUserByEmail(userEmail)
	if err != nil {
		return "", "", fmt.Errorf("service.GetPasswordByEmail: %s", err.Error())
	}
	return returnedUser.UserID, returnedUser.Password, nil
}

func (u *user) List() ([]*UserInfo, error) {
	userList, err := u.userRepo.List()
	if err != nil {
		return nil, fmt.Errorf("service.List: %s", err.Error())
	}
	// fmt.Println("length of userList ---", len(userList))
	length := len(userList)
	var UserList []*UserInfo
	for i := 0; i < length; i++ {
		tempUser := &UserInfo{
			UserID:    userList[i].UserID,
			UserName:  userList[i].UserName,
			UserEmail: userList[i].UserEmail,
		}
		UserList = append(UserList, tempUser)
		// fmt.Println("UserList by now ---- ", fmt.Sprintf("%+v", UserList))
	}
	return UserList, nil
}

func (u *user) UpdateName(userID string, userName string) error {
	returnedUser, err1 := u.userRepo.GetUserByID(userID)
	if err1 != nil {
		return fmt.Errorf("service.UpdateName.GetUserByID: %s", err1.Error())
	} //else {
	// fmt.Println("...#Service response: ID received - ", userID, ", user name received - ", userName)
	// }
	tempUser := &repo.UserInfo{
		UserID:   returnedUser.UserID,
		UserName: userName,
		Password: returnedUser.Password,
	}
	// fmt.Println("...#Service response: tempUser created - ", tempUser)

	err := u.userRepo.Update(*tempUser)
	if err != nil {
		return fmt.Errorf("service.UpdateName.Update: %s", err.Error())
	} //  else {
	// 	fmt.Println("...#Service response: tempUser updated - ", tempUser.UserName)
	// }
	return err
}

func (u *user) UpdatePassword(userID string, password string) error {
	returnedUser, err1 := u.userRepo.GetUserByID(userID)
	if err1 != nil {
		return fmt.Errorf("service.UpdatePassword.GetUserByID: %s", err1.Error())
	}
	tempUser := &repo.UserInfo{
		UserID:   returnedUser.UserID,
		UserName: returnedUser.UserName,
		Password: password,
	}
	err := u.userRepo.Update(*tempUser)
	if err != nil {
		return fmt.Errorf("service.UpdatePassword.Update: %s", err.Error())
	}
	return err
}

func (u *user) DeleteUser(userID string) (string, error) {
	returnedUser, err1 := u.userRepo.GetUserByID(userID)
	if err1 != nil {
		return "", fmt.Errorf("service.DeleteUser.GetUserByID: %s", err1.Error())
	}
	err := u.userRepo.DeleteUser(userID)
	if err != nil {
		return "", fmt.Errorf("service.DeleteUser.Delete: %s", err.Error())
	}
	return returnedUser.UserName, err
}

func (t *tool) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	// fmt.Println(string(bytes))
	// fmt.Println(bytes)
	return string(bytes), err
}

// func convertUserInfo(userInfo repo.UserInfo) *UserInfo {
// 	return &UserInfo{
// 		UserID:   userInfo.UserID,
// 		UserName: userInfo.UserName,
// 	}
// }
