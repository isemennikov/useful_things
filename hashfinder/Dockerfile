# Используем официальный образ Golang для сборки
FROM golang:1.20-alpine as builder

RUN ulimit -n 65536

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем все файлы в рабочую директорию
COPY . .

# Инициализируем модуль Go и копируем зависимости
RUN go mod init filehash && go mod tidy

# Сборка Go приложения
RUN go build -o filehash .

# Используем минимальный образ для запуска
FROM alpine:latest

# Устанавливаем strace
RUN apk add --no-cache strace

# Создаем пользователя и группу
RUN addgroup -S gods && adduser -S hermes -G gods

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем скомпилированное приложение из предыдущего контейнера
COPY --from=builder /app/filehash /app/filehash

# Меняем пользователя
USER hermes

# Указываем команду для запуска приложения
ENTRYPOINT ["/app/filehash"]
