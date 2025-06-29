# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git (required for some Go modules)
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o planning-poker .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary and web files
COPY --from=builder /app/planning-poker .
COPY --from=builder /app/web ./web

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./planning-poker"]
