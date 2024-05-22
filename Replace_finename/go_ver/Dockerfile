# Указываем аргумент для версии Go
ARG GO_VERSION=1.22-alpine

# Сборка исполняемого файла Go
FROM golang:${GO_VERSION} AS builder

# Указываем аргумент для модуля Go (параметризация)
ARG MODULE_NAME=example.com/myapp

# Создание рабочей директории
WORKDIR /app

# Копирование исходного кода
COPY . .

# Проверка и инициализация модуля, если go.mod отсутствует
RUN [ -f go.mod ] || go mod init ${MODULE_NAME}

# Установка зависимостей
RUN go mod tidy

# Сборка приложения. Убедитесь, что ваш код находится в поддиректории /app
# и main.go находится в корне этой директории
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o renamefileg .

# Этап 2: Создание финального образа
FROM scratch

# Создаем non-root пользователя и группу
USER 1001:1001

# Копирование исполняемого файла из билдера
COPY --from=builder /app/renamefileg /renamefileg

# Установка точки входа
ENTRYPOINT ["/renamefileg"]
