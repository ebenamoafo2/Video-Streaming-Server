# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

# Copy source code
COPY main.go .

# Build the binary
RUN go build -o server .

# Runtime stage (smaller image)
FROM alpine:latest

WORKDIR /app

# Copy compiled binary from builder
COPY --from=builder /app/server .

# Copy static files and videos into the image
COPY static/videos ./static/videos

# Ensure video directory exists with proper permissions
RUN chmod -R 755 ./static/videos

EXPOSE 8080

CMD ["./server"]
