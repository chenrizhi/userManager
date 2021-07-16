package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	Db *sql.DB
)

func InitDb(driver, dsn string) error {
	//dsn := "golang:golang@1golang:golang@tcp(127.0.0.1:3306)/user?charset=utf8mb4&loc=Local&parseTime=true"
	//db, err := sql.Open("mysql", dsn)
	var err error
	Db, err = sql.Open(driver, dsn)
	if err != nil {
		return err
	}
	err = Db.Ping()
	return err
}
