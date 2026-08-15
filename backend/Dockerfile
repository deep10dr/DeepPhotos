# Multi-stage Dockerfile for DeepPhotos Go Backend
FROM golang:alpine AS builder

WORKDIR /app

# Install git and build essentials
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build lightweight binary without CGO (pure Go modernc.org/sqlite)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Production minimal image
FROM alpine:latest  

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/server .

# Expose Go REST API port
EXPOSE 8080

CMD ["./server"]
