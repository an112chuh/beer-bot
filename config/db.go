package config

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func InitDb() {
	db, err := sql.Open("sqlite3", "beer.db")
	if err != nil {
		panic(err)
	}
	Db = db
}

func GetConnection() (db *sql.DB) {
	return Db
}
