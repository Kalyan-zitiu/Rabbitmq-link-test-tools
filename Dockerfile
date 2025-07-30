# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.21 AS builder
WORKDIR /app

# Copy go module files and download dependencies first for caching
COPY go.mod ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the rabbitprobe binary
RUN go build -o rabbitprobe ./cmd/rabbitprobe

# Final stage
FROM debian:stable-slim
WORKDIR /usr/local/bin
COPY --from=builder /app/rabbitprobe .

# Default command shows help
ENTRYPOINT ["rabbitprobe"]
CMD ["--help"]
