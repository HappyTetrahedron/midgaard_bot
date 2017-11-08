package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Session struct {
	Chat *tgbotapi.Chat
	Msgbox chan tgbotapi.Message
}

var sessions map[int64]*Session

func initSessions() {
	sessions = make(map[int64]*Session)
}

func getSession(chat *tgbotapi.Chat) *Session {
	session, ok := sessions[chat.ID]
	if !ok {
		session = newSession(chat)
	}
	return session
}

func newSession(chat *tgbotapi.Chat) *Session {
	session := Session{chat, make(chan tgbotapi.Message)}
	sessions[chat.ID] = &session
	return &session
}
