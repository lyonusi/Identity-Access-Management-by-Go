package repo

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var tokenTableName = "token"

type UserToken struct {
	UserID      string
	TokenString string
}

type Token interface {
	CreateToken(userID string, tokenString string) error
}

type token struct {
	db *sql.DB
}

func NewUserToken(db *sql.DB) Token {
	return &token{
		db: db,
	}
}

func (t *token) CreateToken(userID string, tokenString string) error {
	_, err := t.db.Exec(
		fmt.Sprintf(
			`INSERT INTO %s (userID, token) VALUES (?, ?)`,
			tokenTableName,
		),
		userID,
		tokenString,
	)
	if err != nil {
		// fmt.Println(err.Error())
		return fmt.Errorf("repo.CreateToken: %s", err.Error())
	}
	return nil
}
