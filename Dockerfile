# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o confessly .

# Final stage
FROM alpine:3.19

# Install timezone data
RUN apk --no-cache add tzdata ca-certificates

# Create necessary directories
RUN mkdir -p /app/uploads/videos /app/configs /app/internal/configs

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/confessly .

# Copy config files
COPY --from=builder /app/internal/configs/configs.json /app/internal/configs/

# Make binary executable
RUN chmod +x /app/confessly

# Expose port
EXPOSE 8081

# Command to run the application
CMD ["/app/confessly"]
