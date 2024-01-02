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
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/reiver/go-telnet"
)

type Session struct {
	Chat  *tgbotapi.Chat
	Input chan *tgbotapi.Message
}

var sessions map[int64]*Session
var mercHost string
var sessionsCtx context.Context

func initSessions(host string, ctx context.Context) error {
	sessions = make(map[int64]*Session)
	mercHost = host
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
	ctx, cancel := context.WithCancel(sessionsCtx)

	go func() {
		telnetInput, telnetOutput, telnetError := make(chan string), make(chan string), make(chan string)
		caller := TelnetCaller{
			Input:  telnetInput,
			Output: telnetOutput,
			Error:  telnetError,
		}

		go func() {
			for {
				select {
				case msg := <-session.Input:
					if msg.Text != "/start" {
						telnetInput <- strings.Trim(msg.Text, "/")
					}
				case body := <-telnetOutput:
					sendToTelegram(session.Chat.ID, body)
				case <-telnetError:
					cancel()
					delete(sessions, session.Chat.ID)
					return
				case <-ctx.Done():
					return
				}
			}
		}()

		telnet.DialToAndCall(mercHost, caller)
	}()
}

func sendToSession(session *Session, message *tgbotapi.Message) {
	session.Input <- message
}
