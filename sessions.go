package main

import (
	"context"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Session struct {
	Chat *tgbotapi.Chat
	Input chan *tgbotapi.Message
}

var sessions map[int64]*Session
var sessionsCtx context.Context

func initSessions(ctx context.Context) error {
	sessions = make(map[int64]*Session)
	sessionsCtx = ctx
	return nil
}

func getSession(chat *tgbotapi.Chat) *Session {
	session, ok := sessions[chat.ID]
	if !ok {
		session = newSession(chat)
	}
	return session
}

func newSession(chat *tgbotapi.Chat) *Session {
	session := Session{chat, make(chan *tgbotapi.Message)}
	sessions[chat.ID] = &session
	startSession(&session)
	return &session
}

func startSession(session *Session) {
	ctx, _ := context.WithCancel(sessionsCtx)
	go func() {
		for {
			select {
			case msg := <-session.Input:
				log.Printf("[%s] %s", msg.From.UserName, msg.Text)

				newMsg := tgbotapi.NewMessage(msg.Chat.ID, msg.Text)
				newMsg.ReplyToMessageID = msg.MessageID
				sendToTelegram(newMsg)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func sendToSession(session *Session, message *tgbotapi.Message) {
	session.Input <- message
}
