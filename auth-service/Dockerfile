# Build Stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# DEBUG: List files to ensure we copied everything correctly
RUN echo "--- Listing files in /app/cmd/server ---" && ls -la ./cmd/server && echo "---------------------------"

# Build the binary
# We target the directory containing main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Run Stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Copy the .env file from the build context (your local folder) to the container
COPY .env .

#Expose PORT
EXPOSE 3001

# We rely on docker-compose env_file to inject environment variables
CMD ["./main"]
