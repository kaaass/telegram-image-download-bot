# Build stage
FROM golang:1.20 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o telegram-bot .

# Run stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/telegram-bot .

ENV TELEGRAM_API_TOKEN=""
ENV ALLOWED_CHAT_ID=""
ENV DOWNLOAD_PATH="/downloads"
ENV HTTP_PROXY=""

VOLUME /downloads

CMD ["./telegram-bot"]
