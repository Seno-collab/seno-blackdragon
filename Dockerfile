# -------- Stage 1: Build --------
ARG GO_VERSION=1.23
FROM golang:${GO_VERSION}-alpine AS builder

# Cài tool cần thiết
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /src

# Tối ưu cache cho dependencies
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# Copy source
COPY . .

# (Tuỳ dự án) Generate swagger trước khi build
# Nếu file main ở cmd/main.go, chỉnh path theo repo của bạn
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/main.go -o docs

# Build: link tĩnh, gọn, reproducible
ENV CGO_ENABLED=0
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH \
  go build -trimpath -buildvcs=false \
  -ldflags="-s -w -buildid=" \
  -o /out/app ./cmd

# -------- Stage 2: Runtime (Alpine) --------
FROM alpine:3.20

# Non-root user
RUN addgroup -S app && adduser -S app -G app && \
  apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /out/app ./app
# (Tuỳ chọn) copy docs swagger nếu bạn cần serve
# COPY --from=builder /src/docs ./docs

USER app
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --retries=3 \
  CMD wget -qO- http://127.0.0.1:8080/health || exit 1

ENTRYPOINT ["./app"]
