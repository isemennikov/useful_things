FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod init renamefileg && \
    go build -o renamefileg

ENTRYPOINT ["/app/renamefileg"]
