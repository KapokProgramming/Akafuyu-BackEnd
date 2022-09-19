package config

import (
	"server/util"
)

type Config struct {
	Port string
	Env  string
	db   struct {
		Host     string
		Port     string
		User     string
		Password string
		Dbname   string
	}
}

func LoadConfig() Config {
	var cfg Config

	util.CheckENV()

	port := util.MustGet("PORT")

	if port == "" {
		port = "7700"
	}

	cfg.Port = util.MustGet("PORT")
	cfg.Env = util.MustGet("ENV")
	cfg.db.Host = util.MustGet("DB_HOST")
	cfg.db.Port = util.MustGet("DB_PORT")
	cfg.db.User = util.MustGet("DB_USER")
	cfg.db.Password = util.MustGet("DB_PASSWORD")
	cfg.db.Dbname = util.MustGet("DB")

	return cfg
}
