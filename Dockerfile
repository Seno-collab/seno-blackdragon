# syntax=docker/dockerfile:1.7
ARG GO_VERSION=1.24.1
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /src
RUN apt-get update && apt-get install -y --no-install-recommends git ca-certificates tzdata && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
# cache module (cần BuildKit)
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/main.go -o docs

ENV CGO_ENABLED=0
RUN go build -trimpath -buildvcs=false -ldflags="-s -w -buildid=" -o /out/app ./cmd

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /out/app ./app
EXPOSE 8080
USER nonroot
ENTRYPOINT ["./app"]
