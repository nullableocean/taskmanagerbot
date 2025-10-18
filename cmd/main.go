package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"taskbot/delivery/tg"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	conf := NewConfig()

	bot := setupBotWithWebhook(conf)
	responder := tg.NewResponder(bot)
	updateHandler := tg.NewUpdateHandler(responder)

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
