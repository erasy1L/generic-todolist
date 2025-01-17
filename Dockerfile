FROM golang:1.22.5-alpine AS builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o todo-list .

FROM alpine AS hoster

WORKDIR /app

COPY --from=builder /build/.env ./.env
COPY --from=builder /build/migrations ./migrations
COPY --from=builder /build/todo-list ./todo-list

ENTRYPOINT [ "./todo-list" ]