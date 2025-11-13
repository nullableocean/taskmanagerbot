package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"taskbot/delivery/tg"
	"taskbot/pkg/logger"
	"taskbot/repository/pg"
	"taskbot/repository/rdb"
	"taskbot/service/task"
	"taskbot/service/telegram"
	"taskbot/service/user"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	conf := NewConfig()
	ctx := context.Background()

	// logs
	file, err := logger.SetupFileForLogs(conf.LOG.LogDir, conf.LOG.LogFile)
	if err != nil {
		log.Println(err)
	}
	if file != nil {
		defer file.Close()
	}

	// DB
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		conf.DB.User, conf.DB.Pass, conf.DB.DBName, conf.DB.Host, conf.DB.Port,
	)
	db, err := sql.Open("postgres", dsn)
	defer db.Close()

	if err != nil {
		log.Fatalf("db open error: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("db ping error: %v", err)
	}

	redisDb := redis.NewClient(&redis.Options{
		Addr:     conf.REDIS.Host + ":" + conf.REDIS.Port,
		Password: conf.REDIS.Password,
		DB:       0,
	})

	err = redisDb.Set(ctx, "health", "check", 0).Err()
	if err != nil {
		log.Fatalf("redis connect error: %v", err)
	}

	// services
	userRepo := pg.NewUserRepository(db)
	taskRepo := pg.NewTaskRepository(db)
	stateStore := rdb.NewStateStore(redisDb)

	userService := user.NewUserService(userRepo)
	tgUserService := user.NewTelegramUserService(userService, userRepo)
	taskService := task.NewTaskService(taskRepo)

	bot := setupBotWithWebhook(conf)
	updateProccesor := telegram.NewUpdateProccesor(tgUserService, taskService, stateStore)
	updateHandler := tg.NewUpdateHandler(tg.NewResponder(bot), updateProccesor)

	updates := bot.ListenForWebhook("/")
	updateListener := tg.NewUpdateListener(updates, updateHandler)

	go func() {
		err := updateListener.Listen()
		if err != nil {
			log.Fatalf("listener start error: %v", err)
		}
	}()

	listenServer(conf.APP.Port)
	updateListener.Stop()

	log.Println("service closed...")
}

func listenServer(port string) {
	defMux := http.DefaultServeMux
	defMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{Addr: ":" + port, Handler: defMux}

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	log.Println("server start listen...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}

	log.Println("listen closed...")
}

func setupBotWithWebhook(conf Config) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(conf.BOT.Token)
	if err != nil {
		log.Fatalf("create bot error: %v", err)
	}

	wh, _ := tgbotapi.NewWebhook(conf.BOT.WebhookUrl)
	_, err = bot.Request(wh)
	if err != nil {
		log.Fatalf("request webhook error: %v", err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatalf("request webhook info error: %v", err)
	}

	if info.LastErrorDate != 0 {
		log.Fatalf("telegram callback failed: %s", info.LastErrorMessage)
	}

	return bot
}
