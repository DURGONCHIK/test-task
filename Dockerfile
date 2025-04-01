# Используем официальный образ Go для сборки
FROM golang:1.23.2 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы модулей и устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем бинарник с отключенным CGO для совместимости с Alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Минимальный образ для финального контейнера
FROM alpine:latest

WORKDIR /root/

# Копируем бинарник из билдера
COPY --from=builder /app/main .

# Делаем бинарник исполняемым
RUN chmod +x main

# Копируем файлы переменных окружения
COPY .env .env

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]
