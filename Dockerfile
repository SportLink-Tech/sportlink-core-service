# Stage 1: Build the Go application
FROM golang:1.22-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the go module files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files do not change
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app within the specified subdirectory
RUN  go build -o main ./cmd/sportlink/main.go

# Stage 2: Setup the runtime container
#FROM alpine:latest
FROM golang:1.22-alpine

# Set work directory in the new stage
WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable, which is the server
CMD ["./main"]
#CMD ["CompileDaemon", "-command=./main"]
