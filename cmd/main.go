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

	bot, updates := startBotWithWebhook(conf)

	responder := tg.NewResponder(bot)

	updateHandler := tg.NewUpdateHandler(responder)
	updateListener := tg.NewUpdateListener(updates, updateHandler)

	go func() {
		err := updateListener.Listen()
		if err != nil {
			log.Fatalf("listener start error: %v", err)
		}
	}()

	listen(conf.APP.Port)
	updateListener.Stop()
}

func listen(port string) {
	server := &http.Server{Addr: ":" + port}
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

func startBotWithWebhook(conf Config) (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {
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

	updatesCh := bot.ListenForWebhook("/")

	return bot, updatesCh
}
