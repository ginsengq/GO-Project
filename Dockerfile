# Этап сборки
FROM golang:1.24 AS builder

WORKDIR /app

# Установка зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Установка swag (для генерации документации)
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Генерация swagger документации
RUN swag init -g internal/app/start/start.go -o docs/swagger

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/app ./cmd/app/*

# Финальный этап
FROM alpine:3.18

# Установка зависимостей времени выполнения
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Копируем бинарник и документацию
COPY --from=builder /app/bin/app .
COPY --from=builder /app/docs/swagger ./docs/swagger
COPY --from=builder /app/internal/app/config ./config


# Используем непривилегированного пользователя
RUN adduser -D appuser
USER appuser

EXPOSE 8000

CMD ["./app"]