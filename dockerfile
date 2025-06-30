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

# Use a lightweight image for production
FROM python:3.9-slim

# Set the working directory inside the container
WORKDIR /app

# Install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy the built Go binary and application files
COPY --from=builder /app/stocks-api /app/stocks-api
COPY . .

# Expose the port for the application
EXPOSE 8000
EXPOSE 26257

# Command to run the application
CMD ["./stocks-api"]