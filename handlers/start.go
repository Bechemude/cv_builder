package handlers

import tele "gopkg.in/telebot.v4"

const startMessage = `Привет! Я помогу адаптировать твоё резюме под конкретную вакансию.

Отправь мне одно из:
• PDF-файл с резюме
• Ссылку на вакансию
• Текст описания вакансии`

func (h *Handlers) Start(c tele.Context) error {
	sender := c.Sender()
	if _, err := h.s.User.FindOrCreate(sender.ID, sender.Username, sender.FirstName, sender.LastName); err != nil {
		return err
	}

	return c.Send(startMessage)
}
