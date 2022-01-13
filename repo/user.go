package repo

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var tableName = "users"

type UserInfo struct {
	UserID    string
	UserName  string
	UserEmail string
	Password  string
}

type User interface {
	CreateUser(userID string, userName string, userEmail string, password string) error
	GetUserByID(userID string) (*UserInfo, error)
	GetUserByName(userID string) (*UserInfo, error)
	GetUserByEmail(userID string) (*UserInfo, error)
	List() ([]*UserInfo, error)
	Update(UserInfo) error
	UpdateWithoutPassword(UserInfo) error
	DeleteUser(UserID string) error
}

type user struct {
	db *sql.DB
}

func NewUser(db *sql.DB) User {
	return &user{
		db: db,
	}
}

func (u *user) CreateUser(userID string, userName string, userEmail string, password string) error {
	_, err := u.db.Exec(
		fmt.Sprintf(
			`INSERT INTO %s (userID, name, email, password) VALUES (?, ?, ?,?)`,
			tableName,
		),
		userID,
		userName,
		userEmail,
		password,
	)
	if err != nil {
		// fmt.Println(err.Error())
		return fmt.Errorf("repo.CreateUser: %s", err.Error())
	}
	return nil
}

func (u *user) GetUserByID(userID string) (*UserInfo, error) {
	rows, err := u.db.Query(
		fmt.Sprintf(
			`SELECT userID, name, email, password FROM %s WHERE userID = ?`,
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
		err = rows.Scan(&userResult.UserID, &userResult.UserName, &userResult.UserEmail, &userResult.Password)
		if err != nil {
			return nil, fmt.Errorf("repo.GetUserByID: %s", err.Error())
		}
		// fmt.Println(userResult.UserID)
		// fmt.Println(userResult.UserName)
		// fmt.Println(userResult.UserEmail)
		// fmt.Println(userResult.Password)
		return userResult, nil
	}
	return nil, fmt.Errorf("User Not Found")
}

func (u *user) GetUserByName(userName string) (*UserInfo, error) {
	rows, err := u.db.Query(
		fmt.Sprintf(
			`SELECT userID, name, email, password FROM %s WHERE name = ?`,
			tableName,
		),
		userName)
	if err != nil {
		return nil, fmt.Errorf("repo.GetUserByName: %s", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		// var userResult UserInfo
		userResult := &UserInfo{}
		rows.Scan(&userResult.UserID, &userResult.UserName, &userResult.UserEmail, &userResult.Password)
		// fmt.Println(userResult.UserID)
		// fmt.Println(userResult.UserName)
		// fmt.Println(userResult.Password)
		return userResult, nil
	}
	return nil, fmt.Errorf("User Not Found")
}

func (u *user) GetUserByEmail(userEmail string) (*UserInfo, error) {
	rows, err := u.db.Query(
		fmt.Sprintf(
			`SELECT userID, name, email, password FROM %s WHERE email = ?`,
			tableName,
		),
		userEmail)
	if err != nil {
		return nil, fmt.Errorf("repo.GetUserByEmail: %s", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		// var userResult UserInfo
		userResult := &UserInfo{}
		rows.Scan(&userResult.UserID, &userResult.UserName, &userResult.UserEmail, &userResult.Password)
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
			`SELECT userID, email, name FROM %s`,
			tableName,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("repo.List: %s", err.Error())
	}
	var userList []*UserInfo
	for rows.Next() {
		userResult := &UserInfo{}
		rows.Scan(&userResult.UserID, &userResult.UserEmail, &userResult.UserName)
		// fmt.Println(userResult.UserID)
		// fmt.Println(userResult.UserName)
		userList = append(userList, userResult)
	}
	return userList, nil
}

func (u *user) Update(userInfo UserInfo) error {
	tempUser := userInfo
	_, err := u.db.Exec(
		fmt.Sprintf(
			`UPDATE "%s" SET name = ?, email =?, password = ? WHERE userID = ?`,
			tableName,
		),
		tempUser.UserName,
		tempUser.UserEmail,
		tempUser.Password,
		tempUser.UserID,
	)
	if err != nil {
		return fmt.Errorf("repo.Update: %s", err.Error())
	}
	return err
}

func (u *user) UpdateWithoutPassword(userInfo UserInfo) error {
	tempUser := userInfo
	_, err := u.db.Exec(
		fmt.Sprintf(
			`UPDATE "%s" SET name = ?, email =? WHERE userID = ?`,
			tableName,
		),
		tempUser.UserName,
		tempUser.UserEmail,
		tempUser.UserID,
	)
	if err != nil {
		return fmt.Errorf("repo.UpdateWithoutPassword: %s", err.Error())
	}
	return err
}

func (u *user) DeleteUser(userID string) error {
	_, err := u.db.Exec(
		fmt.Sprintf(
			`DELETE FROM %s where userID=?`,
			tableName,
		),
		userID,
	)
	if err != nil {
		return fmt.Errorf("repo.DeleteUser: %s", err.Error())
	}
	return err
}
