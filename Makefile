.PHONY: build test clean docker-build docker-run help

# Default target
all: build

# Build the application
build:
	@echo "Building reviewer-karma..."
	go build -o bin/reviewer-karma ./cmd/reviewer-karma

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f .karma-data.json

# Deploy to GitHub (create release)
deploy:
	@echo "Deploying to GitHub..."
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Usage: make deploy VERSION=v1.0.0"; \
		exit 1; \
	fi
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)
	@echo "✅ Tag $(VERSION) created and pushed!"
	@echo "📝 Next: Create GitHub release at https://github.com/master-wayne7/reviewer-karma-action/releases"

# Quick release (Windows compatible)
release:
	@echo "Creating release..."
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Usage: make release VERSION=v1.0.0"; \
		exit 1; \
	fi
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)
	@echo "✅ Release $(VERSION) created!"
	@echo "🌐 View at: https://github.com/master-wayne7/reviewer-karma-action/releases/tag/$(VERSION)"

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t reviewer-karma .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run --rm -e GITHUB_TOKEN=your_token_here reviewer-karma

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean         - Clean build artifacts"
	@echo "  deploy        - Deploy to GitHub (VERSION=v1.0.0)"
	@echo "  release       - Quick release (VERSION=v1.0.0)"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  deps          - Install dependencies"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  help          - Show this help" 