package middleware

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func dbconnection() {
	db, err := sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/testdb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	insert, err := db.Query("INSERT INTO testdb VALUES(2, 'TEST)")
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
}
