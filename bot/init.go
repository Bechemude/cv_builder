package bot

import (
	"cvbuilder/config"
	"cvbuilder/handlers"
	"log"
	"time"

	tele "gopkg.in/telebot.v4"
)

type Bot struct {
	b *tele.Bot
}

func Init(c *config.Config) *Bot {
	pref := tele.Settings{
		Token:  c.TelegramBotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return &Bot{
		b,
	}
}

func (b *Bot) RegisterHandlers(h *handlers.Handlers) {
	b.b.Handle("/hello", h.SayHello)
}

func (b *Bot) Run() {
	b.b.Start()
}
