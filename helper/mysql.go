package helper

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var (
	user = "root"
	pass = ""
	host = "127.0.0.1:3306"
	db   = "weetbot"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pass, host, db))

	if err != nil {
		return nil, err
	}

	return db, nil
}
