package repo

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var scopeTableName = "user_scope"

type UserScope struct {
	UserID string
	Scope  []string
}

type userScope struct {
	userID string
	scope  string
}

type Scope interface {
	SetScopeByID(userID string, scope string) error
	GetScopeByID(userID string) (*UserScope, error)
	CheckUserScope(userID string, userScope string) (bool, error)
	ListUserIDByScope(scope string) ([]string, error)
	DeleteScopeByID(userID string, scope string) error
}

type scope struct {
	db *sql.DB
}

func NewUserScope(db *sql.DB) Scope {
	return &scope{
		db: db,
	}
}

func (s *scope) SetScopeByID(userID string, scope string) error {
	_, err := s.db.Exec(
		fmt.Sprintf(
			`INSERT INTO %s (userID, scope) VALUES (?, ?)`,
			scopeTableName,
		),
		userID,
		scope,
	)
	if err != nil {
		// fmt.Println(err.Error())
		return fmt.Errorf("repo.SetScopeByID: %s", err.Error())
	}
	return nil
}

func (s *scope) CheckUserScope(userID string, scope string) (bool, error) {
	rows, err := s.db.Query(
		fmt.Sprintf(
			`SELECT scope FROM %s WHERE userId = ?`,
			scopeTableName,
		),
		userID,
	)
	if err != nil {
		return false, fmt.Errorf("repo.CheckUserScope: %s", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		userScopeRow := &userScope{}
		rows.Scan(&userScopeRow.scope)
		if userScopeRow.scope == scope {
			// fmt.Println(userID, "-- scope --", userScopeRow.scope, " [exists]")
			return true, nil
		}
	}
	return false, nil
}

func (s *scope) GetScopeByID(userID string) (*UserScope, error) {
	rows, err := s.db.Query(
		fmt.Sprintf(
			`SELECT scope FROM %s WHERE userId = ?`,
			scopeTableName,
		),
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("repo.GetScopeByID: %s", err.Error())
	}
	var scopeArray []string
	defer rows.Close()
	for rows.Next() {
		userScopeRow := &userScope{}
		rows.Scan(&userScopeRow.scope)
		// fmt.Println("-- scope --", userScopeRow.scope)
		scopeArray = append(scopeArray, userScopeRow.scope)
	}
	userScopeResult := UserScope{
		UserID: userID,
		Scope:  scopeArray,
	}

	return &userScopeResult, nil
}

func (s *scope) ListUserIDByScope(scope string) ([]string, error) {
	rows, err := s.db.Query(
		fmt.Sprintf(
			`SELECT userId FROM %s WHERE scope = ?`,
			scopeTableName,
		),
		scope,
	)
	if err != nil {
		return nil, fmt.Errorf("repo.ListUserIdByScope: %s", err.Error())
	}
	var userIDList []string
	defer rows.Close()
	for rows.Next() {
		userIdRow := &userScope{}
		rows.Scan(&userIdRow.userID)
		// fmt.Println(userIdRow.userID, "-- scope --", userIdRow.scope)
		userIDList = append(userIDList, userIdRow.userID)
	}

	return userIDList, nil
}

func (s *scope) DeleteScopeByID(userID string, userScope string) error {
	result, err := s.db.Exec(
		fmt.Sprintf(
			`DELETE FROM %s WHERE (userId = ? AND scope = ? )`,
			scopeTableName,
		),
		userID, userScope,
	)
	if err != nil {
		return fmt.Errorf("repo.DeleteScopeByID: %s", err.Error())
	}
	deletedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repo.DeleteScopeByID: deletedRows: %s", err.Error())
	}
	if deletedRows == 0 {
		return fmt.Errorf("repo.DeleteScopeByID: user scope does not exist")
	}
	return err
}
