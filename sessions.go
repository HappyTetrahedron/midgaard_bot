package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Session struct {
	Chat *tgbotapi.Chat
	Msgbox chan tgbotapi.Message
}
