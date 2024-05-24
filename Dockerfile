# Use the official Go image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project into the container
COPY . .
COPY configs/ca.crt /app/configs/ca.crt
RUN chmod 644 /app/configs/ca.crt
# Explicitly enable Go modules
ENV GO111MODULE=on

# Build the Go application
RUN go build -o main cmd/main.go

# Expose the port the application runs on
EXPOSE 50051
EXPOSE 9000
# Command to run the Go application
CMD ["/app/main"]