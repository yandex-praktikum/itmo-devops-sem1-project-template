# Этап, на котором выполняется сборка приложения
FROM golang:1.23.3-alpine as builder
WORKDIR /build
COPY . .
RUN go build -o /main src/main/main.go
# Финальный этап, копируем собранное приложение
FROM alpine:3
COPY --from=builder main /bin/main
ENTRYPOINT ["/bin/main"]