package main

import (
	"context"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jessevdk/go-flags"
)

var config struct {
	Token string `short:"t" long:"token" description:"Telegram API Token" required:"true"`
}

func recv_task(bot *tgbotapi.BotAPI, sendChannel chan tgbotapi.Chattable, ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			sendChannel <- msg
		case <-ctx.Done():
			return
		}
	}
}

func send_task(bot *tgbotapi.BotAPI, sendChannel chan tgbotapi.Chattable, ctx context.Context) {
	for {
		select {
		case msg := <-sendChannel:
			bot.Send(msg)
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	_, err := flags.Parse(&config)
	if err != nil {
		log.Panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	ctx, cancel := context.WithCancel(context.Background())
	sendChannel := make(chan tgbotapi.Chattable)
	go send_task(bot, sendChannel, ctx)
	recv_task(bot, sendChannel, ctx)
	cancel()
}
