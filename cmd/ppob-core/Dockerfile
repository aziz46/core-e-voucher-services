# Build stage
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go mod and sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ppob-core ./cmd/ppob-core

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/ppob-core .

EXPOSE 8080

CMD ["./ppob-core"]
