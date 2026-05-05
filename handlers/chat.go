package handlers

import (
	"cvbuilder/models"
	"cvbuilder/services"
	"fmt"
	"strconv"
	"strings"

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

	jobText := fmt.Sprintf(
		"✅ Вакансия сохранена\n\n*%s* — %s\n_%s_\n\nТребуемые навыки:\n%s",
		job.Title,
		job.CompanyName,
		job.Seniority,
		formatList(job.SkillsRequired),
	)
	progress(jobText)

	// Предлагаем выбрать язык адаптированного резюме
	jobIDStr := strconv.FormatUint(uint64(job.ID), 10)
	markup := &tele.ReplyMarkup{}
	markup.InlineKeyboard = [][]tele.InlineButton{{
		{Text: "🇷🇺 Русский", Unique: "cv_lang", Data: jobIDStr + ":ru"},
		{Text: "🇬🇧 English", Unique: "cv_lang", Data: jobIDStr + ":en"},
		{Text: "🌐 Язык вакансии", Unique: "cv_lang", Data: jobIDStr + ":auto"},
	}}
	c.Bot().Send(c.Recipient(), "✍️ На каком языке адаптировать резюме?", markup)

	return nil
}

func formatVariantMessage(v *models.CVVariant) string {
	var sb strings.Builder

	title, matchLabel, changesLabel := variantLabels(v.Language)

	sb.WriteString(fmt.Sprintf("✍️ *%s*\n\n", title))
	sb.WriteString(fmt.Sprintf("🎯 *%s:* %d%%\n\n", matchLabel, v.MatchScore))

	if len(v.KeyChanges) > 0 {
		sb.WriteString(fmt.Sprintf("📝 *%s:*\n", changesLabel))
		for _, change := range v.KeyChanges {
			sb.WriteString("• " + change + "\n")
		}
	}

	return sb.String()
}

func variantLabels(language string) (title, matchLabel, changesLabel string) {
	switch language {
	case "ru":
		return "Резюме адаптировано под вакансию", "Совпадение", "Что изменилось"
	default:
		return "CV tailored for the vacancy", "Match score", "What changed"
	}
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
