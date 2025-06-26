# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o elasticdump .

# Final stage
FROM alpine:latest

LABEL maintainer="Rahul Sahoo"

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/elasticdump .

# Make it executable
RUN chmod +x ./elasticdump

ENTRYPOINT ["./elasticdump"]
