#!/bin/bash

# Reviewer Karma Action Deployment Script
# This script automates the release process for GitHub Marketplace

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if version is provided
if [ -z "$1" ]; then
    print_error "Please provide a version number (e.g., ./scripts/deploy.sh v1.0.0)"
    exit 1
fi

VERSION=$1

print_status "Starting deployment for version: $VERSION"

# Check if we're on main branch
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "main" ]; then
    print_warning "You're not on the main branch. Current branch: $CURRENT_BRANCH"
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_error "Deployment cancelled"
        exit 1
    fi
fi

# Check if there are uncommitted changes
if [ -n "$(git status --porcelain)" ]; then
    print_error "You have uncommitted changes. Please commit or stash them first."
    git status --short
    exit 1
fi

# Run tests
print_status "Running tests..."
if ! go test -v ./...; then
    print_error "Tests failed. Please fix them before deploying."
    exit 1
fi
print_success "All tests passed!"

# Build the application
print_status "Building application..."
if ! go build ./cmd/reviewer-karma; then
    print_error "Build failed. Please fix the issues before deploying."
    exit 1
fi
print_success "Build successful!"

# Create and push tag
print_status "Creating tag: $VERSION"
if git tag -l | grep -q "^$VERSION$"; then
    print_error "Tag $VERSION already exists. Please use a different version."
    exit 1
fi

git tag -a "$VERSION" -m "Release $VERSION"
git push origin "$VERSION"
print_success "Tag created and pushed: $VERSION"

# Create GitHub release
print_status "Creating GitHub release..."
if command -v gh &> /dev/null; then
    gh release create "$VERSION" \
        --title "Release $VERSION" \
        --notes "
## 🚀 Release $VERSION

### Features
- 🎯 Track reviewer engagement and generate karma-based leaderboards
- ⚙️ Fully customizable scoring system
- 🔄 Support for incremental updates
- 🤖 Automatic bot filtering
- 📊 Beautiful markdown leaderboard generation

### Usage
\`\`\`yaml
- name: Run Reviewer Karma Action
  uses: master-wayne7/reviewer-karma-action@$VERSION
  with:
    github-token: \${{ secrets.GITHUB_TOKEN }}
\`\`\`

### Changes
- Initial release with all core features
- Customizable scoring system
- Incremental update support
- Comprehensive documentation
"
    print_success "GitHub release created!"
else
    print_warning "GitHub CLI not found. Please create the release manually:"
    echo "1. Go to: https://github.com/master-wayne7/reviewer-karma-action/releases"
    echo "2. Click 'Create a new release'"
    echo "3. Select tag: $VERSION"
    echo "4. Add release notes"
    echo "5. Publish release"
fi

print_success "Deployment completed successfully!"
print_status "Next steps:"
echo "1. Wait for GitHub Actions to build and release"
echo "2. Submit to GitHub Marketplace: https://github.com/marketplace/actions"
echo "3. Update marketplace listing with new version"
echo "4. Share your action on social media!"

print_success "🎉 Your action is ready for GitHub Marketplace!" 