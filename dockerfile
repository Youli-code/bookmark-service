FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY main.go .

RUN go build -o bookmark-service

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/bookmark-service .

EXPOSE 8080

CMD ["./bookmark-service"]