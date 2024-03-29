/*
midgaard_bot, a Telegram bot which sets a bridge to Midgaard Merc MUD
Copyright (C) 2017 by Javier Sancho Fernandez <jsf at jsancho dot org>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var sendChannel chan tgbotapi.Chattable

func initTelegramWorkers(token string, ctx context.Context) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	sendChannel = make(chan tgbotapi.Chattable)
	go sendWorker(bot, sendChannel, ctx)
	go recvWorker(bot, ctx)

	return nil
}

func recvWorker(bot *tgbotapi.BotAPI, ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			if update.Message != nil && update.Message.Text != "" {
				session := getSession(update.Message.Chat)
				sendToSession(session, update.Message)
			}
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

func sendToTelegram(chat_id int64, body string) {
	newMsg := tgbotapi.NewMessage(chat_id, formatStringForSending(body))
	newMsg.ParseMode = tgbotapi.ModeMarkdownV2
	sendChannel <- newMsg
}

func formatStringForSending(raw string) string {
	return fmt.Sprintf("```\n%s\n```", raw)
}
