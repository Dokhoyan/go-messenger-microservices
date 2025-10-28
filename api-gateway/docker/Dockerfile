FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Финальный образ
FROM alpine:latest

WORKDIR /root/

# Копируем бинарник
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]

