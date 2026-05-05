package handlers

import (
	"cvbuilder/services"
	"fmt"

	tele "gopkg.in/telebot.v4"
)

func (h *Handlers) ChatCompletion(c tele.Context) error {
	user, err := h.s.User.GetByTelegramID(c.Sender().ID)
	if err != nil {
		return c.Send("Сначала напиши /start.")
	}

	cvs, err := h.s.CV.ListByUserID(user.ID)
	if err != nil || len(cvs) == 0 {
		return c.Send("Сначала отправь своё резюме в формате PDF.")
	}

	msg, err := c.Bot().Send(c.Recipient(), "⏳ Обрабатываю...")
	if err != nil {
		return err
	}

	progress := func(text string) {
		c.Bot().Edit(msg, text)
	}

	job, err := h.s.WebReader.Process(c.Text(), user.ID, progress)
	if err != nil {
		progress("❌ Не удалось обработать вакансию.")
		return fmt.Errorf("ошибка обработки: %w", err)
	}

	text := fmt.Sprintf(
		"✅ Вакансия сохранена (ID: %d)\n\n*%s* — %s\n%s\n\nТребуемые навыки: %s",
		job.ID,
		job.Title,
		job.CompanyName,
		job.Seniority,
		formatList(job.SkillsRequired),
	)

	progress(text)
	return nil
}

func formatList(items []string) string {
	if len(items) == 0 {
		return "—"
	}
	result := ""
	for _, item := range items {
		result += "• " + item + "\n"
	}
	return result
}

// Убедимся что ProgressFunc доступен из services
var _ = services.ProgressFunc(nil)
