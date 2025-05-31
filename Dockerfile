# Build stage
FROM golang:1.21-alpine AS builder

# Install git (required for go mod download with private repos)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build arguments for version information
ARG VERSION=dev
ARG COMMIT=unknown
ARG DATE=unknown

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X main.Version=${VERSION} -X main.GitCommit=${COMMIT} -X main.BuildDate=${DATE}" \
    -a -installsuffix cgo \
    -o gingest \
    cmd/gingest/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests and git for cloning repos
RUN apk --no-cache add ca-certificates git

# Create non-root user
RUN addgroup -g 1001 -S gingest && \
    adduser -u 1001 -S gingest -G gingest

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/gingest .

# Change ownership to non-root user
RUN chown gingest:gingest /app/gingest

# Switch to non-root user
USER gingest

# Create volume for output
VOLUME ["/output"]

# Set default working directory for processing
WORKDIR /workspace

# Default command
ENTRYPOINT ["/app/gingest"]
CMD ["--help"] 