# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.24.1
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /src
RUN apt-get update && apt-get install -y --no-install-recommends \
  git ca-certificates tzdata \
  && rm -rf /var/lib/apt/lists/*

# Copy go.mod/go.sum trước để cache
COPY go.mod go.sum ./

# Cache modules (BuildKit)
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod download

# Copy source
COPY . .

# Pin version swag để tránh vỡ build do latest
ARG SWAG_VERSION=v1.16.3
RUN go install github.com/swaggo/swag/cmd/swag@${SWAG_VERSION}
RUN swag init -g cmd/main.go -o docs

# Build (NHỚ có RUN, và giữ \ ở cuối dòng trước)
ENV CGO_ENABLED=0
ARG APP_NAME=app
ARG APP_VERSION=dev
ARG GIT_COMMIT=none
ARG BUILD_TIME=unknown

RUN --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/go/pkg/mod \
  go build -trimpath -buildvcs=false \
  -ldflags="-s -w -buildid= \
  -X 'your/module/internal/version.Version=${APP_VERSION}' \
  -X 'your/module/internal/version.Commit=${GIT_COMMIT}' \
  -X 'your/module/internal/version.BuildTime=${BUILD_TIME}'" \
  -o /out/${APP_NAME} ./cmd

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
ARG APP_NAME=app
COPY --from=builder /out/${APP_NAME} ./app
EXPOSE 8080
USER nonroot
ENTRYPOINT ["./app"]
