package handlers

import (
	"fmt"

	tele "gopkg.in/telebot.v4"
)

func (h *Handlers) ChatCompletion(c tele.Context) error {
	resp, err := h.ex.LLM.ChatCompletion(c.Text())
	if err != nil {
		return err
	}

	fmt.Println(resp)

	return c.Send(resp)
}
