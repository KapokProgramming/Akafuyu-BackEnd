package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func createConnectionToDatabase() *sql.DB {
	if db == nil {
		var err error
		db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("DBUSER"), os.Getenv("DBPASSWORD"), os.Getenv("DBHOST"), os.Getenv("DBPORT"), os.Getenv("DB")))
		if err != nil {
			panic(err)
		}
		db.SetMaxOpenConns(20)
		db.SetMaxIdleConns(20)
		db.SetConnMaxLifetime(time.Minute * 5)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		panic(pingErr)
	}
	return db
}
