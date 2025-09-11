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


# LƯU Ý: tất cả trên MỘT dòng hoặc nối bằng '&&' để không bị tách instruction
RUN go build -trimpath -buildvcs=false \
  -ldflags="-s -w -buildid= 
  -o /out/app ./cmd


FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /out/app ./app
EXPOSE 8080
USER nonroot
ENTRYPOINT ["./app"]
