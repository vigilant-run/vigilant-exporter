FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY server/main.go .
RUN go mod init server && \
    go mod tidy && \
    go build -o server .

FROM golang:1.24-alpine
RUN apk --no-cache add ca-certificates curl
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8000
CMD ["./server"]