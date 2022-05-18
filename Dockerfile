FROM golang:1.18.2-alpine3.15 AS build

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the project files
COPY . .

# Build project
RUN go build -v ./cmd/main.go

# Start from fresh image
FROM alpine:3.15

# Set working directory
WORKDIR /app

# Copy binary from build stage
COPY --from=build /app/main /app/main

# Set default value for GRCP_PORT
ENV GRPC_PORT=4000

# Expose grpc server port
EXPOSE ${GRPC_PORT}

# Start
CMD ["./main"]