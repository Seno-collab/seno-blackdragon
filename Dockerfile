# Stage 1: Build
FROM golang:1.24.1-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/app .
CMD ["./app"]