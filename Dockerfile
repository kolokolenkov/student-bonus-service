# Этап сборки
FROM golang:1.27 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o /app/student-bonus-service \
    ./cmd/api

# Финальный образ
FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/student-bonus-service .

EXPOSE 8080

CMD ["./student-bonus-service"]