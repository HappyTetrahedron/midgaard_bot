package main

import (
	"context"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var sendChannel chan tgbotapi.Chattable

func initTelegramWorkers(token string, ctx context.Context) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	sendChannel = make(chan tgbotapi.Chattable)
	go sendWorker(bot, sendChannel, ctx)
	go recvWorker(bot, ctx)

	return nil
}
	
func recvWorker(bot *tgbotapi.BotAPI, ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			session := getSession(update.Message.Chat)
			sendToSession(session, update.Message)
		case <-ctx.Done():
			return
		}
	}
}

func sendWorker(bot *tgbotapi.BotAPI, sendChannel chan tgbotapi.Chattable, ctx context.Context) {
	for {
		select {
		case msg := <-sendChannel:
			bot.Send(msg)
		case <-ctx.Done():
			return
		}
	}
}

func sendToTelegram(message tgbotapi.Chattable) {
	sendChannel <- message
}
