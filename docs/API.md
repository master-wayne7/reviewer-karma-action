# API Documentation

This document describes the internal API of the Reviewer Karma Action.

## Packages

### `internal/config`

Configuration management for the application.

#### Types

```go
type Config struct {
    ReviewPoint              int
    PositiveEmojiPoint       int
    ConstructiveCommentPoint int
}
```

#### Functions

```go
func Load() Config
```
Loads configuration from environment variables with sensible defaults.

### `internal/karma`

Core karma scoring and leaderboard generation logic.

#### Types

```go
type Reviewer struct {
    Username string `json:"username"`
    Points   int    `json:"points"`
}

type Leaderboard struct {
    Reviewers []Reviewer `json:"reviewers"`
}
```

#### Functions

```go
func IsBot(username string) bool
```
Checks if a username belongs to a bot account.

```go
func HasPositiveEmoji(text string) bool
```
Checks if text contains positive emojis that award bonus points.

```go
func IsConstructiveComment(text string) bool
```
Checks if a comment is constructive (>10 words after filtering).

```go
func GenerateLeaderboard(reviewerKarma map[string]int) Leaderboard
```
Creates a sorted leaderboard from reviewer karma points.

```go
func WriteLeaderboardFile(leaderboard Leaderboard) error
```
Writes the leaderboard to `REVIEWERS.md` file.

### `internal/githubapi`

GitHub API interactions for fetching repository data.

#### Functions

```go
func FetchAllPullRequests(ctx context.Context, client *github.Client, owner, repo string) ([]*github.PullRequest, error)
```
Fetches all pull requests from a repository.

```go
func FetchPullRequestReviews(ctx context.Context, client *github.Client, owner, repo string, prNumber int) ([]*github.PullRequestReview, error)
```
Fetches all reviews for a specific pull request.

```go
func FetchPullRequestComments(ctx context.Context, client *github.Client, owner, repo string, prNumber int) ([]*github.PullRequestComment, error)
```
Fetches all comments for a specific pull request.

## Scoring System

### Default Points

- **Review Point**: 1 point for submitting any review
- **Positive Emoji Point**: 2 points for including positive emojis
- **Constructive Comment Point**: 1 point for constructive comments (>10 words)

### Positive Emojis

The following emojis are considered positive and award bonus points:
- 👍 🔥 😄 🎉 🚀 💯 ✅ ⭐ ❤️ 👏

### Bot Detection

The following patterns are used to identify bot accounts:
- `[bot]` suffix
- `-bot` or `bot-` in username
- Common bot usernames like `github-actions[bot]`, `dependabot[bot]`

### Constructive Comment Detection

A comment is considered constructive if it:
1. Contains more than 10 words after filtering
2. Excludes common non-constructive phrases like "LGTM", "looks good", etc.
3. Provides meaningful feedback or discussion

## Error Handling

All functions return appropriate errors that should be handled by the caller. Common error scenarios include:

- Network errors when fetching from GitHub API
- Rate limiting issues
- Invalid repository information
- File system errors when writing leaderboard

## Performance Considerations

- Uses pagination to handle repositories with many PRs
- Processes reviews and comments in batches
- Efficient string operations for text analysis
- Minimal memory allocation for large datasets 