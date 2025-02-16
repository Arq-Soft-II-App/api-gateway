FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache docker-cli

COPY --from=builder /app/main .
COPY .env .

COPY update_nginx.sh .
RUN chmod +x update_nginx.sh

EXPOSE 8000

CMD ["./main"]
