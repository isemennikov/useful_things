# syntax=docker/dockerfile:1.7.1
FROM alpine:latest

# Установка зависимостей, включая git и procps
RUN apk update && \
    apk add --no-cache bash curl openssl coreutils procps git

# Клонирование репозитория testssl.sh
RUN git clone --depth 1 https://github.com/drwetter/testssl.sh.git /usr/local/testssl.sh

# Установка переменной окружения для testssl.sh
ENV TESTSSL_INSTALL_DIR /usr/local/testssl.sh

# Делаем скрипт исполняемым
RUN chmod +x /usr/local/testssl.sh/testssl.sh

# Установка ENTRYPOINT
ENTRYPOINT ["/usr/local/testssl.sh/testssl.sh"]

# Установка CMD по умолчанию, если аргументы не предоставлены
CMD ["--help"]