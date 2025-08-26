# Build stage
FROM golang:1.21-alpine AS builder

# Set build arguments
ARG VERSION=dev
ARG BUILD_TIME

# Install git and ca-certificates (needed for private repos and SSL)
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" \
    -a -installsuffix cgo \
    -o bmc-cli .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S bmc && \
    adduser -u 1001 -S bmc -G bmc

WORKDIR /home/bmc

# Copy the binary from builder stage
COPY --from=builder /app/bmc-cli .

# Copy sample config
COPY --from=builder /app/config.yaml ./config.yaml.sample

# Change ownership to non-root user
RUN chown -R bmc:bmc /home/bmc

# Switch to non-root user
USER bmc

# Set default command
ENTRYPOINT ["./bmc-cli"]
CMD ["--help"]

# Add labels for better maintainability
LABEL maintainer="BMC CLI Team"
LABEL description="A CLI tool for managing BMCs (iLO and iDRAC)"
LABEL version="${VERSION}" 