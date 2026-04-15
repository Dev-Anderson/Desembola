# Build stage
FROM golang:alpine AS builder

WORKDIR /app

# Copy source code
COPY . .

# Install dependencies and sync go.sum
RUN go mod tidy

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
