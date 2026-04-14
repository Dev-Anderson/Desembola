# Build stage
FROM golang:alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /api cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /api /app/api

# Copy sql migrations folder
COPY --from=builder /app/sql /app/sql

# Copy any existing CSV files so they can be processed if necessary
COPY --from=builder /app/*.csv /app/*.CSV /app/

# Expose API port
EXPOSE 8080

# Command to run the API
CMD ["/app/api"]
