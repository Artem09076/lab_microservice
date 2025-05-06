FROM golang:1.24.1-alpine3.21 AS builder

ENV GOCACHE=/root/.cache/go-build

WORKDIR /app

COPY . .

RUN --mount=type=cache,target="/root/.cache/go-build" \
    go build -o /app/connect-service ./services/connect-service/main.go

FROM alpine:3.21

WORKDIR /app

# Копируем бинарник из builder-этапа
COPY --from=builder /app/connect-service ./

EXPOSE 10000

# Запускаем бинарник (теперь он лежит в /app/connect-service)
CMD ["./connect-service"]