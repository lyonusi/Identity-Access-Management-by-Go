package repo

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var tableName = "users"

type UserInfo struct {
	UserID   string
	UserName string
	Password string
}

type User interface {
	CreateUser(userID string, userName string, password string) error
	GetUserByID(userID string) (*UserInfo, error)
	List() ([]*UserInfo, error)
	Update(UserInfo) error
	// Update2(userID string, field string, updateInfo string) error
}

type user struct {
	db *sql.DB
}

func NewUser(db *sql.DB) User {

	return &user{
		db: db,
	}
}

func (u *user) CreateUser(userID string, userName string, password string) error {
	_, err := u.db.Exec(
		fmt.Sprintf(
			`INSERT INTO %s (userID, name, password) VALUES (?, ?, ?)`,
			tableName,
		),
		userID,
		userName,
		password,
	)
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("repo.CreateUser: %s", err.Error())
	}
	return nil
}

func (u *user) GetUserByID(userID string) (*UserInfo, error) {
	rows, err := u.db.Query(
		fmt.Sprintf(
			`SELECT userID, name, password FROM %s WHERE userID = ?`,
			tableName,
		),
		userID)
	if err != nil {
		return nil, fmt.Errorf("repo.GetUserByID: %s", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		// var userResult UserInfo
		userResult := &UserInfo{}
		rows.Scan(&userResult.UserID, &userResult.UserName, &userResult.Password)
		// fmt.Println(userResult.UserID)
		// fmt.Println(userResult.UserName)
		// fmt.Println(userResult.Password)
		return userResult, nil
	}

	return nil, fmt.Errorf("User Not Found")
}

func (u *user) List() ([]*UserInfo, error) {
	rows, err := u.db.Query(
		fmt.Sprintf(
			`SELECT userID, name FROM %s`,
			tableName,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("repo.List: %s", err.Error())
	}
	var userList []*UserInfo
	for rows.Next() {
		userResult := &UserInfo{}
		rows.Scan(&userResult.UserID, &userResult.UserName)
		// fmt.Println(userResult.UserID)
		// fmt.Println(userResult.UserName)
		userList = append(userList, userResult)
	}
	return userList, nil
}

func (u *user) Update(userInfo UserInfo) error {
	tempUser := userInfo
	fmt.Println("...#Repo response: Ready to update - ", tempUser)
	updateResult, err := u.db.Exec(
		`UPDATE "users" SET name = "test1", password="test2" WHERE userID ="1cb5c0c8-189a-4ec1-8ed0-c5ba4135b8e8"`,
	)
	// updateResult, err := u.db.Exec(
	// 	fmt.Sprintf(
	// 		`UPDATE "%s" SET name = ?, password = ? WHERE userID = ?`,
	// 		tableName,
	// 	),
	// 	tempUser.UserName,
	// 	tempUser.Password,
	// 	tempUser.UserID,
	// )
	if err != nil {
		return fmt.Errorf("repo.Update: %s", err.Error())
	}
	fmt.Println("...#Repo response: Updated - ", updateResult)
	return err
}

// func (u *user) Update2(userID string, field string, updateInfo string) error {
// 	fmt.Println("...#Repo response: Ready to update - ", userID, field, updateInfo)
// 	updateResult, err := u.db.Exec(
// 		fmt.Sprintf(
// 			`UPDATE %s SET %s = ? WHERE userID = ?`,
// 			tableName,
// 			field,
// 		),
// 		updateInfo,
// 	)
// 	if err != nil {
// 		return fmt.Errorf("repo.Update: %s", err.Error())
// 	}
// 	fmt.Println("...#Repo response: Updated - ", updateResult)
// 	return err
// }
