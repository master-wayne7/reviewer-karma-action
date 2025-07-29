# Single stage build for simplicity
FROM golang:1.24-alpine

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o reviewer-karma ./cmd/reviewer-karma

# Verify the binary was created
RUN ls -la reviewer-karma

# Make the binary executable
RUN chmod +x reviewer-karma

# Copy binary to a location that GitHub Actions can find
RUN cp reviewer-karma /usr/local/bin/reviewer-karma

# Set the entrypoint
ENTRYPOINT ["reviewer-karma"] 