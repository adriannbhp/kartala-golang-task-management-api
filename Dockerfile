# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/kartala-api ./cmd/api/main.go

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install CA certificates (needed for NeonDB TLS)
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/kartala-api .
COPY --from=builder /app/docs ./docs

# Cloud Run injects PORT env var; our app checks APP_PORT → PORT → 8080
EXPOSE 8080

CMD ["./kartala-api"]
