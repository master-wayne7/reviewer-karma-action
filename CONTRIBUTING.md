# Contributing to Reviewer Karma Action

Thank you for your interest in contributing to the Reviewer Karma Action! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Code Style](#code-style)
- [Project Structure](#project-structure)

## Code of Conduct

This project is committed to providing a welcoming and inclusive environment for all contributors. Please be respectful and considerate of others.

## Getting Started

1. **Fork the repository**
2. **Clone your fork**
   ```bash
   git clone https://github.com/your-username/reviewer-karma-action.git
   cd reviewer-karma-action
   ```
3. **Set up the development environment** (see [Development Setup](#development-setup))

## Development Setup

### Prerequisites

- Go 1.24 or later
- Docker (optional, for testing the action)
- Git

### Installation

1. **Install dependencies**
   ```bash
   make deps
   ```

2. **Build the application**
   ```bash
   make build
   ```

3. **Run tests**
   ```bash
   make test
   ```

## Making Changes

### 1. Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
```

### 2. Make Your Changes

Follow the [Code Style](#code-style) guidelines and ensure your changes are well-tested.

### 3. Test Your Changes

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Lint code (if you have golangci-lint installed)
make lint
```

### 4. Commit Your Changes

Use conventional commit messages:

```
feat: add new scoring rule for detailed comments
fix: handle edge case in bot detection
docs: update README with new configuration options
test: add tests for new karma calculation logic
```

## Testing

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run tests for a specific package
go test -v ./internal/karma

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Writing Tests

- Write tests for all new functionality
- Aim for at least 80% code coverage
- Use descriptive test names
- Test both success and failure cases

### Example Test

```go
func TestNewFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected bool
    }{
        {"valid input", "test", true},
        {"empty input", "", false},
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            result := NewFeature(test.input)
            if result != test.expected {
                t.Errorf("NewFeature(%q) = %v, expected %v", 
                    test.input, result, test.expected)
            }
        })
    }
}
```

## Submitting Changes

### 1. Push Your Changes

```bash
git push origin feature/your-feature-name
```

### 2. Create a Pull Request

1. Go to the original repository on GitHub
2. Click "New Pull Request"
3. Select your fork and feature branch
4. Fill out the pull request template
5. Submit the pull request

### 3. Pull Request Guidelines

- **Title**: Use a clear, descriptive title
- **Description**: Explain what the change does and why it's needed
- **Tests**: Ensure all tests pass
- **Documentation**: Update documentation if needed
- **Breaking Changes**: Clearly mark any breaking changes

## Code Style

### Go Code Style

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` to format code
- Use meaningful variable and function names
- Add comments for exported functions and complex logic

### File Organization

```
cmd/reviewer-karma/     # Main application entry point
internal/               # Internal packages
├── config/            # Configuration management
├── githubapi/         # GitHub API interactions
└── karma/             # Karma scoring logic
pkg/                   # Public packages (if any)
docs/                  # Documentation
examples/              # Usage examples
```

### Naming Conventions

- **Packages**: Use lowercase, single-word names
- **Functions**: Use camelCase for private, PascalCase for exported
- **Variables**: Use camelCase
- **Constants**: Use UPPER_SNAKE_CASE
- **Files**: Use snake_case

## Project Structure

```
reviewer-karma-action/
├── cmd/reviewer-karma/     # Main application
├── internal/               # Internal packages
│   ├── config/            # Configuration
│   ├── githubapi/         # GitHub API
│   └── karma/             # Karma logic
├── .github/workflows/     # GitHub Actions
├── docs/                  # Documentation
├── examples/              # Usage examples
├── Makefile              # Build automation
├── Dockerfile            # Container definition
├── action.yml            # GitHub Action definition
├── README.md             # Project documentation
├── CONTRIBUTING.md       # This file
├── LICENSE               # License file
└── go.mod               # Go module definition
```

## Common Tasks

### Adding a New Scoring Rule

1. **Update the karma package**
   ```go
   // internal/karma/karma.go
   func NewScoringRule(text string) bool {
       // Implementation
   }
   ```

2. **Add tests**
   ```go
   // internal/karma/karma_test.go
   func TestNewScoringRule(t *testing.T) {
       // Tests
   }
   ```

3. **Update configuration**
   ```go
   // internal/config/config.go
   type Config struct {
       NewRulePoint int
   }
   ```

4. **Update documentation**
   - Update README.md
   - Update action.yml inputs

### Adding a New GitHub API Integration

1. **Create new function in githubapi package**
   ```go
   // internal/githubapi/githubapi.go
   func FetchNewData(ctx context.Context, client *github.Client, owner, repo string) ([]*NewType, error) {
       // Implementation
   }
   ```

2. **Add tests**
   ```go
   // internal/githubapi/githubapi_test.go
   func TestFetchNewData(t *testing.T) {
       // Tests
   }
   ```

## Getting Help

- **Issues**: Use GitHub Issues for bug reports and feature requests
- **Discussions**: Use GitHub Discussions for questions and general discussion
- **Documentation**: Check the README.md and inline code comments

## License

By contributing to this project, you agree that your contributions will be licensed under the same license as the project.

---

Thank you for contributing to Reviewer Karma Action! 🏆 