# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:3.20

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk --no-cache add ca-certificates tzdata


WORKDIR /root/

# Copy the binary and config from builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/config ./config

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./main"]