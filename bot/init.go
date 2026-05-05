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
	b.b.Handle("/start", h.Start)
	b.b.Handle(tele.OnDocument, h.Document)
	b.b.Handle(tele.OnText, h.ChatCompletion)
	b.b.Handle(&tele.InlineButton{Unique: "cv_lang"}, h.CVLangCallback)
}

func (b *Bot) Run() {
	log.Println("Connected to tg bot")

	b.b.Start()
}
