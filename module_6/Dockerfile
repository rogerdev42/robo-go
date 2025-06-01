# Используем официальный Golang образ как базовый
FROM golang:1.23-alpine AS builder

# Устанавливаем необходимые пакеты
RUN apk add --no-cache git ca-certificates tzdata

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go mod файлы
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# Используем минимальный образ для production
FROM alpine:latest

# Добавляем ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Копируем собранное приложение
COPY --from=builder /app/main .

# Копируем миграции
COPY --from=builder /app/internal/database/migrations ./migrations

# Уст