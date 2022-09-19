package config

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func OpenDB(cfg Config) (*sql.DB, error) {

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.db.User, cfg.db.Password, cfg.db.Host, cfg.db.Port, cfg.db.Dbname))

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}
