package main

import "os"

const (
	postgresPort = "5432"
	redisPort    = "6379"
	port         = "8000"
	logDir       = "./logs"
)

type Config struct {
	DB struct {
		User   string
		Pass   string
		Host   string
		Port   string
		DBName string
	}

	REDIS struct {
		Host     string
		Port     string
		Password string
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
	cnf.DB.Port = postgresPort
	cnf.DB.DBName = os.Getenv("DB_NAME")

	cnf.REDIS.Host = os.Getenv("REDIS_HOST")
	cnf.REDIS.Port = os.Getenv("REDIS_PORT")
	cnf.REDIS.Password = os.Getenv("REDIS_PASSWORD")

	cnf.BOT.Token = os.Getenv("BOT_TOKEN")
	cnf.BOT.WebhookUrl = os.Getenv("BOT_WEBHOOK")

	cnf.APP.Port = port

	cnf.LOG.LogDir = logDir
	cnf.LOG.LogFile = os.Getenv("LOG_FILE")

	return cnf
}
