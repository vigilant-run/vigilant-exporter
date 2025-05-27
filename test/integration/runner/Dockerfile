FROM golang:1.24-alpine

WORKDIR /app
COPY runner/integration_test.go .

RUN go mod init tests && \
    go mod tidy
