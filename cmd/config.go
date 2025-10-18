package main

import "os"

const (
	port   = "8000"
	logDir = "./logs"
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

	LOG struct {
		LogFile string
		LogDir  string
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

	cnf.LOG.LogDir = logDir
	cnf.LOG.LogFile = os.Getenv("LOG_FILE")

	return cnf
}
