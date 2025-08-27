# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.24.1
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /src
RUN apt-get update && apt-get install -y --no-install-recommends \
  git ca-certificates tzdata \
  && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SWAG_VERSION=v1.16.3
RUN go install github.com/swaggo/swag/cmd/swag@${SWAG_VERSION}
RUN swag init -g cmd/main.go -o docs

ENV CGO_ENABLED=0
ARG APP_NAME=app
ARG APP_VERSION=dev
ARG GIT_COMMIT=none
ARG BUILD_TIME=unknown

# LƯU Ý: tất cả trên MỘT dòng hoặc nối bằng '&&' để không bị tách instruction
RUN go build -trimpath -buildvcs=false \
  -ldflags="-s -w -buildid= \
  -X 'seno-blackdragon/internal/version.Version=${APP_VERSION}' \
  -X 'seno-blackdragon/internal/version.Commit=${GIT_COMMIT}' \
  -X 'seno-blackdragon/internal/version.BuildTime=${BUILD_TIME}'" \
  -o /out/${APP_NAME} ./cmd


FROM gcr.io/distroless/static:nonroot
WORKDIR /app
ARG APP_NAME=app
COPY --from=builder /out/${APP_NAME} ./app
EXPOSE 8080
USER nonroot
ENTRYPOINT ["./app"]
