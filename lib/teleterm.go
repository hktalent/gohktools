package main

import (
	"log"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"
)
/*
https://github.com/tucnak/telebot
*/
func main() {
	pref := tele.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
b.Handle(tele.OnText, func(c tele.Context) error {
	// All the text messages that weren't
	// captured by existing handlers.

	var (
		user = c.Sender()
		text = c.Text()
	)

	// Use full-fledged bot's functions
	// only if you need a result:
	msg, err := b.Send(user, text)
	if err != nil {
		return err
	}

	// Instead, prefer a context short-hand:
	return c.Send(text)
})

b.Handle(tele.OnChannelPost, func(c tele.Context) error {
	// Channel posts only.
	msg := c.Message()
})

b.Handle(tele.OnPhoto, func(c tele.Context) error {
	// Photos only.
	photo := c.Message().Photo
})

b.Handle(tele.OnQuery, func(c tele.Context) error {
	// Incoming inline queries.
	return c.Answer(...)
})
	b.Handle("/hello", func(c tele.Context) error {
		return c.Send("Hello!")
	})

	b.Start()
}

