# Build stage
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Install air for hot reload
RUN go install github.com/air-verse/air@latest

# Copy the source code
COPY . .

# Run air for development (will be overridden in production)
CMD ["air"]
