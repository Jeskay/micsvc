FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GO_ARCH=amd64 go build -o notify_service ./cmd/notify_server/main.go

FROM scratch

COPY --from=builder /app/notify_service ./

ENV PORT=9091 \
    HOST=localhost \
    TOKEN_EXPIRE=12 \
    CONNECTION_TIMEOUT=2 \
    KAFKA_ADDRESS=localhost:9093

EXPOSE $PORT

CMD ["./notify_service"]
