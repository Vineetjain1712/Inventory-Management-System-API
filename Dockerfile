# Step 1: Build the Go binary
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy dependency files first and download them
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your project files
COPY . .

# Build the Go app (no CGO needed with modernc.org/sqlite)
RUN CGO_ENABLED=0 GOOS=linux go build -o inventory-api ./cmd/main.go

# Step 2: Run the app in a lightweight container
FROM alpine:latest

# Create a working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/inventory-api .

# Expose the port your app listens on
EXPOSE 8080

# Run the app
CMD ["./inventory-api"]
