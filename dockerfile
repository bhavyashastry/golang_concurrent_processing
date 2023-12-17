# Use an official Golang runtime as a base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go server source code into the container
COPY . .

# Build the Go application
RUN go build -o app

# Expose the port on which the Go server listens
EXPOSE 8000

# Command to run the Go server when the container starts
CMD ["./app"]
