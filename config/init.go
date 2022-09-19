package config

import (
	"flag"
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

	// cfg.Port = port
	// cfg.Env = env.MustGet("ENV")

	flag.StringVar(&cfg.Port, "port", port, "Server port will be listen on")
	flag.StringVar(&cfg.Env, "environment", util.MustGet("ENV"), "Application environment")
	// cfg.db.Host = env.MustGet("DBHOST")
	// cfg.db.Port = env.MustGet("DBPORT")
	// cfg.db.User = env.MustGet("DBUSER")
	// cfg.db.Password = env.MustGet("DBPASS")
	// cfg.db.Dbname = env.MustGet("DBNAME")

	flag.StringVar(&cfg.db.Host, "database host", util.MustGet("DBHOST"), "some host")
	flag.StringVar(&cfg.db.Port, "database port", util.MustGet("DBPORT"), "some port")
	flag.StringVar(&cfg.db.User, "database user", util.MustGet("DBUSER"), "user")
	flag.StringVar(&cfg.db.Password, "database password", util.MustGet("DBPASS"), "password")
	flag.StringVar(&cfg.db.Dbname, "database name", util.MustGet("DBNAME"), "db name")
	return cfg
}
