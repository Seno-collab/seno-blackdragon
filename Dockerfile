# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.24.1

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-bookworm AS deps
WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

FROM deps AS code
COPY . .

ARG SWAG_VERSION=v1.16.3

RUN --mount=type=cache,target=/go/pkg/mod \
    go install github.com/swaggo/swag/cmd/swag@${SWAG_VERSION} && \
    swag init -g cmd/main.go -o docs

FROM code AS build
ARG TARGETOS TARGETARCH
ENV CGO_ENABLED=0 \
    GOFLAGS="-mod=readonly"

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build \
      -trimpath \
      -buildvcs=false \
      -ldflags "-s -w -buildid=" \
      -o /out/app ./cmd

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app

COPY --from=build /out/app ./app
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["./app"]
