package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func Initialize(dsn string) (func() error, error) {
	var err error
	DB, err = sqlx.Open("mysql", dsn)
	if err != nil {
		return func() error { return nil }, err
	}
	return DB.Close, DB.Ping()
}

func IsNoRows(err error) bool {
	return err == sql.ErrNoRows
}
