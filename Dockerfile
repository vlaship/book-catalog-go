### Build stage
# Start from a small, secure base image
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy only necessary files for dependency resolution
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Install tools and generate docs/code
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    go install github.com/google/wire/cmd/wire@latest && \
    swag init -g ./internal/router/router.go -o ./api/docs && \
    wire ./...

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -ldflags="-w -s" -o /app/app ./cmd/app/main.go

### Run stage
# Create a minimal production image
FROM alpine:3.19

# Update packages and clean up in a single layer
RUN apk update && \
    apk upgrade && \
    apk add --no-cache ca-certificates tzdata && \
    rm -rf /var/cache/apk/* /tmp/*

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/app .

# Copy API docs
COPY --from=builder /app/api/docs ./api/docs

# Create a non-root user
RUN adduser -D appuser
USER appuser

# Run the binary when the container starts
CMD [ "./app" ]
