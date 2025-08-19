# Stage 1: Build
FROM golang:1.24.1-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/app ./cmd
RUN swag init -g cmd/main.go -o docs

FROM alpine:3.20

WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]