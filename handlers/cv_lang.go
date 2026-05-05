package handlers

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v4"
)

// CVLangCallback handles the language selection inline button.
// Callback data format: "{jobID}:{lang}", e.g. "42:ru", "42:en", "42:auto"
func (h *Handlers) CVLangCallback(c tele.Context) error {
	_ = c.Respond()

	parts := strings.SplitN(c.Data(), ":", 2)
	if len(parts) != 2 {
		return c.Edit("❌ Некорректный запрос.")
	}

	jobID64, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return c.Edit("❌ Некорректный ID вакансии.")
	}
	jobID := uint(jobID64)
	language := parts[1]

	user, err := h.s.User.GetByTelegramID(c.Sender().ID)
	if err != nil {
		return c.Edit("❌ Пользователь не найден. Напиши /start.")
	}

	if err := c.Edit("⏳ Генерирую адаптированное резюме..."); err != nil {
		return err
	}

	job, err := h.s.Job.GetByID(jobID)
	if err != nil {
		return c.Edit("❌ Вакансия не найдена.")
	}

	cvs, err := h.s.CV.ListByUserID(user.ID)
	if err != nil || len(cvs) == 0 {
		return c.Edit("❌ Резюме не найдено. Сначала загрузи PDF.")
	}

	cv := cvs[len(cvs)-1]
	if err := h.s.CV.LoadJobsHistory(&cv); err != nil {
		return c.Edit("❌ Не удалось загрузить историю работ.")
	}

	variant, err := h.s.CVVariant.Generate(&cv, job, user.ID, language)
	if err != nil {
		_ = c.Edit("❌ Не удалось адаптировать резюме.")
		return fmt.Errorf("generate variant error: %w", err)
	}

	_ = c.Edit("📄 Рендерю PDF...")

	fullCV := variant.BuildCV(&cv)
	pdfBytes, err := h.s.PDFGenerator.Render(fullCV, language)
	if err != nil {
		_ = c.Edit("❌ Не удалось сгенерировать PDF.")
		return fmt.Errorf("pdf render error: %w", err)
	}

	langLabel := langName(language)
	_ = c.Edit(fmt.Sprintf("✅ Резюме адаптировано на %s", langLabel))

	// Send summary card
	c.Bot().Send(c.Recipient(), formatVariantMessage(variant), tele.ModeMarkdown)

	// Send PDF file
	fileName := fmt.Sprintf("cv_%s_%s.pdf", fullCV.FirstName, fullCV.LastName)
	doc := &tele.Document{
		File:     tele.FromReader(bytes.NewReader(pdfBytes)),
		FileName: fileName,
		MIME:     "application/pdf",
	}
	c.Bot().Send(c.Recipient(), doc)

	return nil
}

func langName(code string) string {
	switch code {
	case "ru":
		return "русском языке"
	case "en":
		return "английском языке"
	default:
		return "языке вакансии"
	}
}
