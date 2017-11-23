package main

import (
	"context"
	"log"

	"github.com/jessevdk/go-flags"
)

var config struct {
	Token string `short:"t" long:"token" description:"Telegram API Token" required:"true"`
}

func main() {
	_, err := flags.Parse(&config)
	if err != nil {
		log.Panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = initSessions(ctx)
	if err != nil {
		log.Panic(err)
	}

	err = initTelegramWorkers(config.Token, ctx)
	if err != nil {
		log.Panic(err)
	}

	for {
		select {
		case <-ctx.Done():
			break
		}
	}
}
