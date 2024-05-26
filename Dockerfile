### Build stage
# Start from a small, secure base image
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Download the Go module dependencies
RUN go mod download

# Install Swag and generate API docs
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init -g ./internal/router/router.go -o ./api/docs

# Install and Run wire
RUN go install github.com/google/wire/cmd/wire@latest
RUN wire ./...

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o /app/app ./cmd/app/main.go

### Run stage
# Create a minimal production image
FROM alpine:3.19

# It's essential to regularly update the packages within the image to include security patches
RUN apk update && apk upgrade && \
# Reduce image size
    rm -rf /var/cache/apk/* && \
    rm -rf /tmp/*

# Set the working directory inside the container
WORKDIR /app

# Copy only the necessary files from the builder stage
COPY --from=builder /app/app .

# Avoid running code as a root user
RUN adduser -D appuser
USER appuser

# Run the binary when the container starts
CMD [ "./app" ]
