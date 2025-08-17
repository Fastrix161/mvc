# Use official Golang image as the builder
FROM golang:1.24.5 AS builder

# Set working directory inside the container
WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app (adjust the path to main.go if needed)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main ./cmd/main.go


# Use a minimal final image to run the binary
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/pkg/views ./pkg/views
COPY --from=builder /app/public/images ./public/images

# Expose the app's port (change if needed)
EXPOSE 8100

# Command to run the app
CMD ["./main"]
