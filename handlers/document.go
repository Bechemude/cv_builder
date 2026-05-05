package handlers

import (
	"fmt"
	"io"

	tele "gopkg.in/telebot.v4"
)

func (h *Handlers) Document(c tele.Context) error {
	doc := c.Message().Document
	if doc.MIME != "application/pdf" {
		return c.Send("Пожалуйста, отправь PDF-файл.")
	}

	user, err := h.s.User.GetByTelegramID(c.Sender().ID)
	if err != nil {
		return c.Send("Сначала напиши /start.")
	}

	msg, err := c.Bot().Send(c.Recipient(), "⏳ Загружаю файл...")
	if err != nil {
		return err
	}

	progress := func(text string) {
		c.Bot().Edit(msg, text)
	}

	progress("📥 Скачиваю файл из Telegram...")
	reader, err := c.Bot().File(&doc.File)
	if err != nil {
		progress("❌ Не удалось загрузить файл.")
		return fmt.Errorf("ошибка загрузки файла: %w", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		progress("❌ Не удалось прочитать файл.")
		return fmt.Errorf("ошибка чтения файла: %w", err)
	}

	cv, err := h.s.PDFReader.Read(doc.FileName, data, user.ID, progress)
	if err != nil {
		progress("❌ Ошибка при обработке резюме.")
		return fmt.Errorf("ошибка обработки PDF: %w", err)
	}

	progress(fmt.Sprintf("✅ Резюме сохранено (ID: %d). Найдено позиций: %d.", cv.ID, len(cv.JobsHistory)))
	return nil
}
