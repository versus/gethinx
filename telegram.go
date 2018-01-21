package main

import (
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

// https://github.com/go-telegram-bot-api/telegram-bot-api
// https://medium.com/golang-notes/%D0%BF%D0%B8%D1%88%D0%B5%D0%BC-%D0%B1%D0%BE%D1%82%D0%B0-%D0%B4%D0%BB%D1%8F-telegram-%D0%BD%D0%B0-go-71c9acd102d1

func StartTelegramBot() {

	// TODO: придумать вывод таблицы статистики
	// TODO: придумать авторизацию
	// TODO: добавить команды reload  и update  инфы с серверов

	b, err := tb.NewBot(tb.Settings{
		Token:  conf.Telegram.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Println("Error create telegram bot ", err.Error())
		return
	}

	b.Handle("/status", func(m *tb.Message) {
		b.Send(m.Sender, GetStatusTable())
	})

	b.Start()
}
