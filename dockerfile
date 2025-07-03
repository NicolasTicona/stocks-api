# Build Stage
FROM golang:1.24.4-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o stocks-api .

# Expose the port for the application
EXPOSE 8000

# Command to run the application
CMD ["./stocks-api"]