# Reviewer Karma Action 🏆

A GitHub Action that tracks reviewer engagement and generates a karma-based leaderboard for your repository.

[![GitHub Actions](https://img.shields.io/badge/GitHub%20Actions-Ready-blue?logo=github-actions)](https://github.com/master-wayne7/reviewer-karma-action)
[![Go](https://img.shields.io/badge/Go-1.24+-blue?logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![GitHub Marketplace](https://img.shields.io/badge/GitHub%20Marketplace-Available-orange?logo=github)](https://github.com/marketplace/actions/reviewer-karma-action)

> **Track, score, and celebrate your repository's most engaged reviewers!** 🎯

## Features

- 🎯 **Automatic Tracking**: Monitors all pull request reviews and comments
- 🏅 **Karma Scoring**: Awards points for reviews, positive emojis, and constructive feedback
- 🤖 **Bot Filtering**: Automatically ignores bot comments and reviews
- 📊 **Leaderboard Generation**: Creates a beautiful markdown leaderboard
- ⚙️ **Configurable Points**: Customize scoring via environment variables
- 🚀 **Efficient**: Handles repositories with 100+ PRs efficiently

## Scoring System

| Action | Default Points | Description | Customizable |
|--------|----------------|-------------|--------------|
| ✅ Code Review | +1 | Awarded for any review submission | ✅ Yes |
| 🎉 Positive Emoji | +2 | For reviews/comments with 👍, 🔥, 😄, etc. | ✅ Yes |
| 💬 Constructive Comment | +1 | For comments with >10 meaningful words | ✅ Yes |

**All points are fully customizable!** See [Configuration](#configuration) section below.

## Quick Start

### 1. Add the Action to Your Repository

Create `.github/workflows/reviewer-karma.yml`:

#### **Option A: Using from GitHub Marketplace (Recommended)**
```yaml
name: Reviewer Karma Tracker

on:
  pull_request_review:
    types: [submitted, edited, dismissed]
  issue_comment:
    types: [created, edited, deleted]
  pull_request:
    types: [opened, synchronize, reopened, closed]
  workflow_dispatch: # Allow manual triggering

jobs:
  update-leaderboard:
    runs-on: ubuntu-latest
    
    steps:
    - name: Checkout repository
      uses: actions/checkout@v4
      with:
        token: ${{ secrets.GITHUB_TOKEN }}

    - name: Run Reviewer Karma Action
      uses: master-wayne7/reviewer-karma-action@v1
      with:
        github-token: ${{ secrets.GITHUB_TOKEN }}
        review-point: '1'
        positive-emoji-point: '2'
        constructive-comment-point: '1'

    - name: Commit and push changes
      run: |
        git config --local user.email "action@github.com"
        git config --local user.name "GitHub Action"
        git add REVIEWERS.md
        git diff --quiet && git diff --staged --quiet || git commit -m "Update reviewer karma leaderboard"
        git push
```

#### **Option B: Using from Local Repository (Development)**
```yaml
    - name: Run Reviewer Karma Action
      uses: ./
      with:
        github-token: ${{ secrets.GITHUB_TOKEN }}
        review-point: '1'
        positive-emoji-point: '2'
        constructive-comment-point: '1'
```

### 2. Configure Repository Permissions

Ensure your workflow has the necessary permissions to:
- Read repository contents
- Write to the repository (for updating `REVIEWERS.md`)

### 3. The Action Will Generate

A `REVIEWERS.md` file in your repository root with content like:

```markdown
# Reviewer Karma Leaderboard

This leaderboard tracks reviewer engagement and contributions to the repository.

## Scoring System

- ✅ Giving a code review: +1 point
- ✅ Review includes a positive emoji (👍, 🔥, 😄, etc.): +2 points
- ✅ Review comment contains a constructive message (>10 words): +1 point

## Current Rankings

| Rank | Reviewer | Points |
|------|----------|--------|
| 1 | 🥇 @alice | 18 |
| 2 | 🥈 @bob | 12 |
| 3 | 🥉 @carol | 10 |
| 4 | @dave | 8 |
| 5 | @eve | 5 |

---
*Last updated: 2024-01-15 14:30:25 UTC*
```

## Configuration

### Environment Variables

You can customize the scoring system using these environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `REVIEW_POINT` | `1` | Points for submitting a review |
| `POSITIVE_EMOJI_POINT` | `2` | Points for including positive emojis |
| `CONSTRUCTIVE_COMMENT_POINT` | `1` | Points for constructive comments |
| `INCREMENTAL_UPDATE` | `false` | Use incremental updates (only process new PRs) |

### Action Inputs

When using the action, you can also configure via inputs:

```yaml
- name: Run Reviewer Karma Action
  uses: ./
  with:
    github-token: ${{ secrets.GITHUB_TOKEN }}
    review-point: '2'           # Custom review points
    positive-emoji-point: '3'   # Custom emoji points
    constructive-comment-point: '2' # Custom comment points
    incremental-update: 'true'  # Enable incremental updates
```

## Positive Emojis

The action recognizes these positive emojis for bonus points:
- 👍 🔥 😄 🎉 🚀 💯 ✅ ⭐ ❤️ 👏

## Bot Detection

The action automatically filters out comments and reviews from:
- Users with `[bot]` suffix
- Common bot usernames (`github-actions[bot]`, `dependabot[bot]`, etc.)
- Users with `-bot` or `bot-` in their username

## Constructive Comment Detection

A comment is considered "constructive" if it:
- Contains more than 10 words after filtering
- Excludes common non-constructive phrases like "LGTM", "looks good", etc.
- Provides meaningful feedback or discussion

## Development

### Prerequisites

- Go 1.24+
- Docker (for building the action)

### Building Locally

```bash
# Install dependencies
go mod download

# Build the application
go build -o reviewer-karma .

# Run locally (requires GITHUB_TOKEN environment variable)
export GITHUB_TOKEN=your_token_here
./reviewer-karma
```

### Building Docker Image

```bash
docker build -t reviewer-karma .
docker run -e GITHUB_TOKEN=your_token_here reviewer-karma
```

## Architecture

The action is built with:

- **Go 1.24+**: For efficient, cross-platform execution
- **GitHub REST API**: For fetching PR and review data
- **Docker**: For containerized execution in GitHub Actions
- **OAuth2**: For secure GitHub API authentication

## Performance

- Efficiently handles repositories with 100+ PRs
- Uses pagination to fetch all data
- Processes reviews and comments in batches
- Minimal API rate limit impact

## Update Modes

### Full Recreation (Default)
- Processes **all PRs** every time the action runs
- **Always accurate** - no risk of stale data
- **Handles deletions** - works correctly if PRs are deleted/modified
- **Best for**: Small to medium repositories (<500 PRs)

### Incremental Updates
- Processes **only new PRs** that haven't been seen before
- **Much faster** - skips already processed PRs
- **Uses storage** - maintains `.karma-data.json` file
- **Best for**: Large repositories (500+ PRs)

To enable incremental updates:
```yaml
- name: Run Reviewer Karma Action
  uses: ./
  with:
    incremental-update: 'true'
```

**Note**: When using incremental updates, the action will create a `.karma-data.json` file to track processed PRs. Make sure to commit this file along with `REVIEWERS.md`.

## Custom Scoring Strategies

### 🎯 **Quality-Focused Scoring**
Emphasize detailed feedback over quick approvals:
```yaml
review-point: '1'                    # Low base points
positive-emoji-point: '1'            # Low emoji value
constructive-comment-point: '5'      # High value for detailed feedback
```

### 🚀 **Engagement-Focused Scoring**
Encourage active participation and positive interactions:
```yaml
review-point: '3'                    # Good base points
positive-emoji-point: '3'            # High emoji value
constructive-comment-point: '2'      # Balanced feedback value
```

### ⚖️ **Balanced Scoring**
Equal emphasis on all types of contributions:
```yaml
review-point: '2'                    # Balanced base points
positive-emoji-point: '2'            # Balanced emoji value
constructive-comment-point: '2'      # Balanced feedback value
```

### 🏆 **Competitive Scoring**
High points for all activities to create more competitive leaderboards:
```yaml
review-point: '5'                    # High base points
positive-emoji-point: '3'            # Medium emoji value
constructive-comment-point: '5'      # High feedback value
```

The leaderboard will automatically display your custom scoring values instead of the defaults!

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) for details.

## Support

If you encounter any issues or have questions:

1. Check the [Issues](../../issues) page
2. Create a new issue with detailed information
3. Include your workflow configuration and any error messages

---

Made with ❤️ by the open source community 