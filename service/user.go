package service

import (
	"IAMbyGo/repo"
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserInfo struct {
	UserID    string `json:"userID"`
	UserName  string `json:"userName"`
	UserEmail string `json:"userEmail"`
}

type UserScope struct {
	UserID    string   `json:"userID"`
	UserScope []string `json:"userScope"`
}

type User interface {
	CreateUser(userName string, userEmail string, password string) error
	GetUserByID(userID string) (*UserInfo, error)
	GetUserByEmail(userEmail string) (*UserInfo, error)
	List() ([]*UserInfo, error)
	UpdateUser(userInfo UserInfo) error
	UpdatePassword(userID string, password string) error
	DeleteUser(userID string) (string, error)
	GetPasswordByName(userName string) (userID string, password string, err error)
	GetPasswordByEmail(userEmail string) (userID string, password string, err error)
	SetScopeByID(userID string, userScope string) error
	CheckUserScope(userID string, userScope string) (bool, error)
	ListScopeByID(userID string) (*UserScope, error)
	ListUserByScope(scope string) ([]*UserInfo, error)
	DeleteScopeByID(userID string, userScope string) error
}

type tools interface {
	hashPassword(password string) (string, error)
}

type tool struct {
}

type user struct {
	privateMethods tools
	userRepo       repo.User
	userScope      repo.Scope
	redisClient    *redis.Client
}

func NewUser(userRepo repo.User, redisClient *redis.Client, userScope repo.Scope) User {
	return &user{
		userRepo:       userRepo,
		redisClient:    redisClient,
		userScope:      userScope,
		privateMethods: &tool{},
	}
}

func (u *user) SetScopeByID(userID string, userScope string) error {
	checkCurrentScope, err := u.userScope.CheckUserScope(userID, userScope)
	if err != nil {
		return fmt.Errorf("service.SetScopeByID.GetCurrentScope: %s", err.Error())
	}
	if checkCurrentScope {
		return fmt.Errorf("service.SetScopeByID.GetCurrentScope: scope already exists")
	}
	err = u.userScope.SetScopeByID(userID, userScope)
	if err != nil {
		return fmt.Errorf("service.SetScopeByID: %s", err.Error())
	}
	return err
}

func (u *user) CheckUserScope(userID string, userScope string) (bool, error) {
	checkCurrentScope, err := u.userScope.CheckUserScope(userID, userScope)
	if err != nil {
		return false, fmt.Errorf("service.SetScopeByID.GetCurrentScope: %s", err.Error())
	}
	if checkCurrentScope {
		return true, nil
	}
	return false, nil
}

func (u *user) ListScopeByID(userID string) (*UserScope, error) {
	dbReturnedUserScope, err := u.userScope.GetScopeByID(userID)
	if err != nil {
		return nil, fmt.Errorf("service.GetScopeByID.GetScopeByID: %s", err.Error())
	}
	userScope := &UserScope{
		UserID:    userID,
		UserScope: dbReturnedUserScope.Scope,
	}
	return userScope, nil
}
func (u *user) ListUserByScope(scope string) ([]*UserInfo, error) {
	dbReturnedUserIdList, err := u.userScope.ListUserIDByScope(scope)
	if err != nil {
		return nil, fmt.Errorf("service.ListUserByScope.ListUserIDByScoe: %s", err.Error())
	}
	var userList []*UserInfo
	for i := 0; i < len(dbReturnedUserIdList); i++ {
		userInfo, err := u.GetUserByID(dbReturnedUserIdList[i])
		if err != nil {
			return nil, fmt.Errorf("service.ListUserByScope.GetUserByID[%v]: %s", i, err.Error())
		}
		userList = append(userList, userInfo)
	}
	// fmt.Println("user list = ", userList)
	return userList, nil
}
func (u *user) DeleteScopeByID(userID string, userScope string) error {
	err := u.userScope.DeleteScopeByID(userID, userScope)
	if err != nil {
		return fmt.Errorf("service.DeleteScopeByID: %s", err.Error())
	}
	return nil
}

func (t *tool) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	// fmt.Println(string(bytes))
	// fmt.Println(bytes)
	return string(bytes), err
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

// func (u *user) GetUserByID(userID string) (*UserInfo, error) {
// 	var ctx = context.Background()
// 	redisReturnedUser, err := u.redisClient.Get(ctx, fmt.Sprintf("user-%s", userID)).Result()

// 	ifReadFromDb := false
// 	tempUserJson := UserInfo{}

// 	if err == redis.Nil || err != nil {
// 		ifReadFromDb = true
// 		fmt.Println(userID, " does not exist")
// 	} else {
// 		fmt.Println(userID, "-----", redisReturnedUser)
// 		err = json.Unmarshal([]byte(redisReturnedUser), &tempUserJson)
// 		if err != nil {
// 			ifReadFromDb = true
// 		}
// 	}

// 	if ifReadFromDb {
// 		dbReturnedUser, err := u.userRepo.GetUserByID(userID)
// 		fmt.Printf("%v\n", dbReturnedUser)
// 		if err != nil {
// 			return nil, fmt.Errorf("service.GetUserByID: %s", err.Error())
// 		}

// 		tempUser := &UserInfo{
// 			UserID:    dbReturnedUser.UserID,
// 			UserName:  dbReturnedUser.UserName,
// 			UserEmail: dbReturnedUser.UserEmail,
// 		}

// 		var tempUserString []byte
// 		tempUserString, err = json.Marshal(tempUser)

// 		fmt.Println(tempUserString)

// 		if err != nil {
// 			fmt.Println(fmt.Errorf("service.GetUserByID.toString: %s", err.Error()))
// 			return tempUser, nil
// 		}
// 		_, err = u.redisClient.Set(ctx, fmt.Sprintf("user-%s", userID), tempUserString, 0).Result()

// 		if err != nil {
// 			fmt.Println(fmt.Errorf("service.GetUserByID.setRedis: %s", err.Error()))
// 			return tempUser, nil
// 		}
// 		return tempUser, nil
// 	} else {
// 		return &tempUserJson, nil
// 	}
// }

func (u *user) GetUserByID(userID string) (*UserInfo, error) {
	readDb := func() ([]byte, error) {
		dbReturnedUser, err := u.userRepo.GetUserByID(userID)
		// fmt.Printf("%v\n", dbReturnedUser)
		if err != nil {
			return nil, fmt.Errorf("service.GetUserByID: %s", err.Error())
		}

		dbReturnedUserInfo := &UserInfo{
			UserID:    dbReturnedUser.UserID,
			UserName:  dbReturnedUser.UserName,
			UserEmail: dbReturnedUser.UserEmail,
		}

		userInfoJson, err := json.Marshal(dbReturnedUserInfo)
		if err != nil {
			return nil, err
		}
		return userInfoJson, nil
	}

	result, err := repo.ReadDbWithCache(readDb, fmt.Sprintf("user-%s", userID), u.redisClient)

	if err != nil {
		return nil, err
	}

	userInfo := UserInfo{}

	// fmt.Printf("result = %+v\n", result)
	err = json.Unmarshal(result, &userInfo)
	if err != nil {
		return nil, err
	}
	// fmt.Println("userInfo = ", userInfo)

	return &userInfo, nil
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

func (u *user) UpdateUser(userInfo UserInfo) error {
	tempUser := &repo.UserInfo{
		UserID:    userInfo.UserID,
		UserName:  userInfo.UserName,
		UserEmail: userInfo.UserEmail,
	}
	// fmt.Println("...#Service response: tempUser created - ", tempUser)

	err := u.userRepo.UpdateWithoutPassword(*tempUser)
	if err != nil {
		return fmt.Errorf("service.UpdateUser.Update: %s", err.Error())
	}

	u.redisClient.Del(context.Background(), fmt.Sprintf("user-%s", tempUser.UserID))

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
	u.redisClient.Del(context.Background(), fmt.Sprintf("user-%s", userID))
	return returnedUser.UserName, err
}
