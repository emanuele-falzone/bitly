FROM golang:1.18.2-alpine3.15 AS build

# Install Alpine Dependencies
RUN apk update && apk upgrade && \
    apk add --no-cache make protoc gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the project files
COPY cmd cmd
COPY internal internal
COPY Makefile ./
COPY LICENSE ./
COPY README.md ./

# Generate code
RUN make generate-code

# Generate docs
RUN make generate-docs

# Define development target
FROM build AS development

# Build with -race option
RUN make build-for-development

# Start development
CMD ["./main"]

# Define production target
FROM build AS production_build

# Build project
RUN make build-for-production

# Start from fresh image
FROM alpine:3.15 AS production

# Set working directory
WORKDIR /app

# Copy binary from build stage
COPY --from=production_build /app/main /app/main

# Start
CMD ["./main"]