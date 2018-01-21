package main

import (
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func StartTelegramBot() {
	b, err := tb.NewBot(tb.Settings{
		Token:  conf.Telegram.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "hello world")
	})

	b.Start()
}
