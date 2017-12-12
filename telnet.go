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
	"bytes"
	"time"

	"github.com/reiver/go-oi"
	"github.com/reiver/go-telnet"
)

type TelnetCaller struct {
	Input chan string
	Output chan string
}

func (caller TelnetCaller) CallTELNET(ctx telnet.Context, writer telnet.Writer, reader telnet.Reader) {

	// Send text to MUD
	go func() {
		var buffer bytes.Buffer
		var p []byte

		crlfBuffer := [2]byte{'\r', '\n'}
		crlf := crlfBuffer[:]

		for {
			message := <-caller.Input
			buffer.Write([]byte(message))
			buffer.Write(crlf)

			p = buffer.Bytes()
			oi.LongWrite(writer, p)
			buffer.Reset()
		}
	}()

	// Receive text from MUD
	chunks := make(chan string)
	chunk := ""
	
	go func() {
		var buffer [1]byte
		p := buffer[:]

		for {
			n, err := reader.Read(p)
			if n <= 0 && err == nil {
				continue
			} else if n <= 0 && err != nil {
				break
			}

			chunks <- string(p)
		}
	}()

	for {
		select {
		case input := <-chunks:
			chunk = input
			for chunk != "" {
				select {
				case input := <-chunks:
					chunk = chunk + input
				case <-time.After(time.Millisecond * 500):
					caller.Output <- chunk
					chunk = ""
				}
			}
		}
	}
}
