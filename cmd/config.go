package main

import "os"

const (
	port = "8000"
)

type Config struct {
	DB struct {
		User   string
		Pass   string
		Host   string
		Port   string
		DBName string
	}

	BOT struct {
		Token      string
		WebhookUrl string
	}

	APP struct {
		Port string
	}
}

func NewConfig() Config {
	cnf := Config{}

	cnf.DB.User = os.Getenv("PG_USER")
	cnf.DB.Pass = os.Getenv("PG_PASSWORD")
	cnf.DB.Host = os.Getenv("PG_HOST")
	cnf.DB.Port = os.Getenv("PG_PORT")
	cnf.DB.DBName = os.Getenv("DB_NAME")

	cnf.BOT.Token = os.Getenv("BOT_TOKEN")
	cnf.BOT.WebhookUrl = os.Getenv("BOT_WEBHOOK")

	cnf.APP.Port = port

	return cnf
}
