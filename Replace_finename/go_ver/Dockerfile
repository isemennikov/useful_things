# Этап 1: Сборка исполняемого файла Go
FROM golang:1.22-alpine AS builder

# Создание рабочей директории
WORKDIR /app

# Копирование исходного кода
COPY . .

# Сборка приложения. Убедитесь, что ваш код находится в поддиректории /app
# и main.go находится в корне этой директории
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o renamefileg .

# Этап 2: Создание финального образа
FROM scratch

# Копирование исполняемого файла из билдера
COPY --from=builder /app/renamefileg /renamefileg

# Установка точки входа
ENTRYPOINT ["/renamefileg"]
