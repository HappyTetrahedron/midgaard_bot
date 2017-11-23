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
