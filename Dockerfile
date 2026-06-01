FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/cvbuilder .

FROM alpine:3.21

RUN apk add --no-cache \
    chromium \
    poppler-utils \
    nss \
    freetype \
    freetype-dev \
    harfbuzz \
    ca-certificates \
    tzdata

ENV CHROME_BIN=/usr/bin/chromium-browser \
    CHROME_PATH=/usr/lib/chromium/ \
    NO_SANDBOX=true

WORKDIR /app

COPY --from=builder /app/cvbuilder .
COPY --from=builder /app/prompts ./prompts

CMD ["./cvbuilder"]
