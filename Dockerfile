# # Stage 1: Build
# FROM golang:latest AS builder

# # Set working directory inside container
# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .

# RUN go build -o app .