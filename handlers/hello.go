package handlers

import (
	tele "gopkg.in/telebot.v4"
)

func (h *Handlers) SayHello(c tele.Context) error {
	return c.Send("Hello!")
}
