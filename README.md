# CV Builder — Telegram Bot

A Telegram bot that parses PDF resumes, analyzes job vacancies, and generates tailored CV variants with cover letters.

## Flow

1. **Upload PDF** → bot extracts text via `pdftotext` → LLM parses into structured CV (work history, skills, summary)
2. **Send job URL/text** → bot fetches and analyzes the vacancy → LLM extracts requirements, company info, red flags
3. **Pick language** → bot generates a tailored CV variant + cover letter optimized for that specific role
4. **Receive** → summary card + cover letter (text) + PDF (without cover letter)

## Tech Stack

- **Go 1.26** (telebot.v4, long polling)
- **PostgreSQL 16** (GORM, auto-migrate)
- **OpenRouter API** (LLM for CV parsing, job analysis, variant generation)
- **chromedp** (HTML → PDF rendering)
- **pdftotext** (poppler-utils — PDF text extraction)

## Architecture

```
config → db → external (LLM, web fetcher) → repos → services → handlers → bot
```

Key packages:
- `config/` — loads `.env` via godotenv
- `models/` — GORM models (User, CV, JobHistory, Job, CVVariant)
- `repos/` — thin CRUD wrappers
- `services/` — business logic (PDF reader, web reader, enricher, CV variant generator, PDF renderer)
- `external/` — LLM client (OpenRouter), web fetcher
- `handlers/` — Telegram command handlers
- `bot/` — telebot init + route registration
- `prompts/` — LLM prompt templates (embedded via `go:embed`)

## Setup

1. Copy `.env-example` to `.env` and fill all fields:
   - PostgreSQL connection
   - Telegram bot token
   - OpenRouter API key + model names
2. Ensure dependencies:
   - PostgreSQL 16
   - `pdftotext` (poppler-utils)
   - Chromium (for chromedp PDF rendering)

## Running

### Taskfile (local dev, requires Go Task)

```bash
task setup                             # First-time setup (.env, DB, deps)
task dev                               # Dev mode: native Postgres + air hot reload
task build                             # Build binary to bin/
task test                              # Run tests
task lint                              # go vet + gofmt
task ci                                # Lint + test + build
task docker:up                         # Start app + postgres in Docker
task docker:down                       # Stop Docker stack
```

### Docker Compose

```bash
docker compose up                      # Postgres + app in containers
docker compose up -d postgres && air   # Postgres in Docker, app locally with hot reload
```

## Prompts

Prompts live in `prompts/` and are embedded at build time via `//go:embed`. Edit markdown files, then rebuild.
