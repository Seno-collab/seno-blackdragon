# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.24.1

#############################
# Builder
#############################
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /src

# Cài tool cần thiết cho swag (git, ca-cert, tzdata nếu bạn muốn tạo file zoneinfo nhúng)
RUN apt-get update && apt-get install -y --no-install-recommends \
  git ca-certificates tzdata \
  && rm -rf /var/lib/apt/lists/*

# Copy mod files trước để tối ưu layer cache
COPY go.mod go.sum ./

# Cache module downloads
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod download

# Copy toàn bộ source
COPY . .

# (Khuyến nghị) pin phiên bản swag để reproducible
ARG SWAG_VERSION=v1.16.3
RUN go install github.com/swaggo/swag/cmd/swag@${SWAG_VERSION}

# Generate swagger (đảm bảo đường dẫn main đúng)
RUN swag init -g cmd/main.go -o docs

# Tối ưu build: dùng cache cho go build + tắt CGO
ENV CGO_ENABLED=0
# Nếu cần timezone local mà vẫn dùng distroless:static, có thể build kèm timetzdata:
# ENV GOFLAGS="-tags timetzdata"

go build -trimpath -buildvcs=false \
  -ldflags "-s -w -buildid= \
  -X 'your/module/internal/version.Version=1.0.0' \
  -X 'your/module/internal/version.Commit=$(git rev-parse --short HEAD)' \
  -X 'your/module/internal/version.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)'" \
  -o bin/app ./cmd

#############################
# Runtime
#############################
# Nếu bạn cần timezone local: thay static -> static-debian12
# FROM gcr.io/distroless/static-debian12:nonroot
FROM gcr.io/distroless/static:nonroot

WORKDIR /app

# Nếu cần timezone local mà dùng static-debian12, có thể copy zoneinfo:
# COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy binary
ARG APP_NAME=app
COPY --from=builder /out/${APP_NAME} ./app

# App listen port (đồng bộ với code & compose)
EXPOSE 8080

# Distroless đã là nonroot
USER nonroot

# Không có shell -> entrypoint là binary
ENTRYPOINT ["./app"]
