# Project Structure

This document describes the organization and structure of the Reviewer Karma Action project.

## Directory Structure

```
reviewer-karma-action/
├── cmd/
│   └── reviewer-karma/          # Main application entry point
│       └── main.go
├── internal/                     # Internal packages (not importable)
│   ├── config/                  # Configuration management
│   │   ├── config.go
│   │   └── config_test.go
│   ├── githubapi/               # GitHub API interactions
│   │   └── githubapi.go
│   └── karma/                   # Karma scoring logic
│       ├── karma.go
│       └── karma_test.go
├── .github/
│   └── workflows/               # GitHub Actions workflows
│       └── reviewer-karma.yml
├── docs/                        # Documentation
│   ├── API.md                   # API documentation
│   └── PROJECT_STRUCTURE.md     # This file
├── examples/                    # Usage examples
│   ├── basic-usage.yml
│   └── custom-scoring.yml
├── .gitignore                   # Git ignore rules
├── action.yml                   # GitHub Action definition
├── CONTRIBUTING.md              # Contributing guidelines
├── Dockerfile                   # Container definition
├── LICENSE                      # MIT License
├── Makefile                     # Build automation
├── README.md                    # Project documentation
├── go.mod                       # Go module definition
└── go.sum                       # Go module checksums
```

## Package Organization

### `cmd/reviewer-karma/`

The main application entry point. Contains the `main()` function and orchestrates the entire workflow.

**Responsibilities:**
- Parse environment variables
- Initialize GitHub client
- Coordinate between packages
- Handle errors and exit codes

### `internal/config/`

Configuration management for the application.

**Responsibilities:**
- Load configuration from environment variables
- Provide sensible defaults
- Validate configuration values

**Key Types:**
```go
type Config struct {
    ReviewPoint              int
    PositiveEmojiPoint       int
    ConstructiveCommentPoint int
}
```

### `internal/karma/`

Core karma scoring and leaderboard generation logic.

**Responsibilities:**
- Calculate karma points for reviewers
- Detect bots, positive emojis, and constructive comments
- Generate and sort leaderboards
- Write leaderboard to markdown file

**Key Types:**
```go
type Reviewer struct {
    Username string `json:"username"`
    Points   int    `json:"points"`
}

type Leaderboard struct {
    Reviewers []Reviewer `json:"reviewers"`
}
```

### `internal/githubapi/`

GitHub API interactions for fetching repository data.

**Responsibilities:**
- Fetch pull requests from repositories
- Fetch reviews for specific pull requests
- Fetch comments for specific pull requests
- Handle pagination and rate limiting

## Design Principles

### 1. Separation of Concerns

Each package has a single, well-defined responsibility:
- **Config**: Configuration management
- **Karma**: Business logic for scoring
- **GitHubAPI**: External API interactions
- **Main**: Application orchestration

### 2. Dependency Inversion

The main application depends on abstractions (interfaces) rather than concrete implementations, making it easy to test and modify.

### 3. Error Handling

All functions return appropriate errors that are handled at the appropriate level. The main application handles fatal errors and exits gracefully.

### 4. Testability

Each package has comprehensive tests:
- Unit tests for all exported functions
- Table-driven tests for edge cases
- Mock tests for external dependencies

## File Naming Conventions

- **Packages**: Lowercase, single word (e.g., `config`, `karma`)
- **Files**: Snake case (e.g., `config.go`, `karma_test.go`)
- **Functions**: PascalCase for exported, camelCase for private
- **Variables**: camelCase
- **Constants**: UPPER_SNAKE_CASE

## Build and Development

### Makefile Targets

```bash
make build          # Build the application
make test           # Run all tests
make test-coverage  # Run tests with coverage
make clean          # Clean build artifacts
make docker-build   # Build Docker image
make docker-run     # Run Docker container
make deps           # Install dependencies
make fmt            # Format code
make lint           # Lint code
```

### Development Workflow

1. **Setup**: `make deps`
2. **Build**: `make build`
3. **Test**: `make test`
4. **Format**: `make fmt`
5. **Lint**: `make lint`

## Deployment

### GitHub Action

The action is deployed as a Docker container with:
- Multi-stage build for minimal image size
- Alpine Linux base for security
- Statically linked binary for portability

### Local Development

```bash
# Build locally
make build

# Run with environment variables
GITHUB_TOKEN=your_token GITHUB_REPOSITORY=owner/repo ./bin/reviewer-karma
```

## Testing Strategy

### Unit Tests

- **Coverage**: Aim for >80% code coverage
- **Packages**: Each package has its own test file
- **Style**: Table-driven tests for multiple scenarios

### Integration Tests

- **GitHub API**: Mock responses for testing
- **File System**: Test leaderboard file generation
- **Configuration**: Test environment variable parsing

### Test Organization

```
internal/
├── config/
│   ├── config.go
│   └── config_test.go
├── karma/
│   ├── karma.go
│   └── karma_test.go
└── githubapi/
    ├── githubapi.go
    └── githubapi_test.go
```

## Documentation

### README.md

- Project overview and features
- Quick start guide
- Configuration options
- Usage examples

### CONTRIBUTING.md

- Development setup
- Code style guidelines
- Testing requirements
- Pull request process

### docs/API.md

- Detailed API documentation
- Function signatures
- Type definitions
- Usage examples

## Future Enhancements

### Potential Additions

- **pkg/**: Public packages for reuse
- **cmd/**: Additional command-line tools
- **scripts/**: Build and deployment scripts
- **deploy/**: Deployment configurations
- **docs/examples/**: More detailed examples

### Scalability Considerations

- **Modular Design**: Easy to add new scoring rules
- **Configuration**: Environment-based configuration
- **Performance**: Efficient pagination and caching
- **Monitoring**: Structured logging for observability 