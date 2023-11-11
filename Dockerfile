# syntax=docker/dockerfile:1

# Build image
FROM golang:1.21 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o ./subscriber

# Deployment image
FROM ubuntu:22.04
WORKDIR /app
COPY --from=build /app/subscriber ./subscriber
ENTRYPOINT ["./subscriber"]