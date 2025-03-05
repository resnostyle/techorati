# Stage 1: Build the Go binary
FROM golang:latest AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Stage 2: Create a minimal runtime image
FROM scratch

# Copy the binary from the builder stage into the minimal runtime image
COPY --from=builder /app/app /app/app

# Expose the port on which the application will listen
EXPOSE 8080

# Command to run the web server
ENTRYPOINT ["/app/app"]
