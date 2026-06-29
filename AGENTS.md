# Run
- `docker compose up` — starts PostgreSQL + app in Docker
- `docker compose up -d postgres && air` - start PostgreSQL in Docker and app in local
- `ls -a **/*.go | go run main.go` — local run (needs postgres, pdftotext, chromium)
- Config: copy `.env-example` to `.env`, fill all fields
- No tests or lint config in this repo

# Architecture
Go 1.26 Telegram bot (telebot.v4, long polling). Init order: `config` → `db` (Postgres, GORM auto-migrate) → `external` (OpenRouter LLM, HTTP web fetcher) → `repos` → `services` → `handlers` → `bot`

# Key packages
- `config/` — loads `.env` via godotenv
- `models/` — GORM models (User, CV, JobHistory, Job, CVVariant) + FlexTime (nullable time)
- `repos/` — thin CRUD wrappers over GORM
- `services/` — business logic: PDF reader/extract, PDF render, web reader, enricher, CV variant gen, template builder
- `external/` — LLM client (go-openrouter), web fetcher (raw HTTP)
- `handlers/` — Telegram command handlers (start, document, text, callback)
- `bot/` — telebot init + route registration
- `prompts/` — LLM prompt templates (go:embed, edits need rebuild)

# Local deps
- `pdftotext` (poppler-utils) for PDF text extraction
- Chromium (for chromedp PDF rendering)
- PostgreSQL 16
