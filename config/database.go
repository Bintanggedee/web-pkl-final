package config

import (
	"database/sql"
)

func Connect_DB() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1)/web_ppkl?parseTime=true")

	if err != nil {
		panic(err)
	}
	// log.Println("Database connected")

	return db
}
