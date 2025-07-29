package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v62/github"
	"github.com/master-wayne7/reviewer-karma-action/internal/config"
	"github.com/master-wayne7/reviewer-karma-action/internal/githubapi"
	"github.com/master-wayne7/reviewer-karma-action/internal/karma"
	"github.com/master-wayne7/reviewer-karma-action/internal/storage"
	"golang.org/x/oauth2"
)

func main() {
	// Check for help flag
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		fmt.Println("Reviewer Karma Action")
		fmt.Println("Track reviewer engagement and generate a karma-based leaderboard")
		fmt.Println("")
		fmt.Println("Environment variables:")
		fmt.Println("  GITHUB_TOKEN          - GitHub token for API access")
		fmt.Println("  GITHUB_REPOSITORY     - Repository name (format: owner/repo)")
		fmt.Println("  REVIEW_POINT          - Points for reviews (default: 1)")
		fmt.Println("  POSITIVE_EMOJI_POINT  - Points for emojis (default: 2)")
		fmt.Println("  CONSTRUCTIVE_COMMENT_POINT - Points for comments (default: 1)")
		fmt.Println("  INCREMENTAL_UPDATE    - Use incremental updates (default: false)")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("  ./reviewer-karma [--help]")
		os.Exit(0)
	}

	// Get GitHub token from environment
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Println("❌ GITHUB_TOKEN environment variable is required")
		os.Exit(1)
	}

	// Get repository information from environment
	repoFullName := os.Getenv("GITHUB_REPOSITORY")
	if repoFullName == "" {
		fmt.Println("❌ GITHUB_REPOSITORY environment variable is required")
		os.Exit(1)
	}

	// Parse repository name (format: "owner/repo")
	parts := strings.Split(repoFullName, "/")
	if len(parts) != 2 {
		fmt.Println("❌ Invalid GITHUB_REPOSITORY format. Expected 'owner/repo'")
		os.Exit(1)
	}
	repoOwner := parts[0]
	repoName := parts[1]

	// Load configuration
	cfg := config.Load()

	// Create GitHub client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	fmt.Printf("🔍 Analyzing repository: %s/%s\n", repoOwner, repoName)
	fmt.Printf("📊 Karma configuration: Review=%d, Emoji=%d, Constructive=%d\n",
		cfg.ReviewPoint, cfg.PositiveEmojiPoint, cfg.ConstructiveCommentPoint)
	fmt.Printf("🔄 Update mode: %s\n", getUpdateModeString(cfg.IncrementalUpdate))

	if cfg.IncrementalUpdate {
		runIncrementalUpdate(ctx, client, repoOwner, repoName, cfg)
	} else {
		runFullRecreation(ctx, client, repoOwner, repoName, cfg)
	}

	fmt.Println("✅ Reviewer karma leaderboard generated successfully!")
}

func runFullRecreation(ctx context.Context, client *github.Client, owner, repo string, cfg config.Config) {
	fmt.Println("🔄 Running in full recreation mode...")

	// Fetch all pull requests
	prs, err := githubapi.FetchAllPullRequests(ctx, client, owner, repo)
	if err != nil {
		fmt.Printf("❌ Error fetching pull requests: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📋 Found %d pull requests\n", len(prs))

	// Calculate karma for all reviewers
	reviewerKarma := make(map[string]int)

	for _, pr := range prs {
		fmt.Printf("🔍 Processing PR #%d: %s\n", pr.GetNumber(), pr.GetTitle())

		// Get reviews for this PR
		reviews, err := githubapi.FetchPullRequestReviews(ctx, client, owner, repo, pr.GetNumber())
		if err != nil {
			fmt.Printf("⚠️ Error fetching reviews for PR #%d: %v\n", pr.GetNumber(), err)
			continue
		}

		// Process reviews
		for _, review := range reviews {
			username := review.GetUser().GetLogin()
			if karma.IsBot(username) {
				continue
			}

			// Award points for giving a review
			reviewerKarma[username] += cfg.ReviewPoint

			// Check for positive emojis in review body
			if karma.HasPositiveEmoji(review.GetBody()) {
				reviewerKarma[username] += cfg.PositiveEmojiPoint
				fmt.Printf("  🎉 @%s gets +%d points for positive emoji\n", username, cfg.PositiveEmojiPoint)
			}

			// Check for constructive comments
			if karma.IsConstructiveComment(review.GetBody()) {
				reviewerKarma[username] += cfg.ConstructiveCommentPoint
				fmt.Printf("  💬 @%s gets +%d points for constructive comment\n", username, cfg.ConstructiveCommentPoint)
			}
		}

		// Get comments for this PR
		comments, err := githubapi.FetchPullRequestComments(ctx, client, owner, repo, pr.GetNumber())
		if err != nil {
			fmt.Printf("⚠️ Error fetching comments for PR #%d: %v\n", pr.GetNumber(), err)
			continue
		}

		// Process comments
		for _, comment := range comments {
			username := comment.GetUser().GetLogin()
			if karma.IsBot(username) {
				continue
			}

			// Check for positive emojis in comment body
			if karma.HasPositiveEmoji(comment.GetBody()) {
				reviewerKarma[username] += cfg.PositiveEmojiPoint
				fmt.Printf("  🎉 @%s gets +%d points for positive emoji in comment\n", username, cfg.PositiveEmojiPoint)
			}

			// Check for constructive comments
			if karma.IsConstructiveComment(comment.GetBody()) {
				reviewerKarma[username] += cfg.ConstructiveCommentPoint
				fmt.Printf("  💬 @%s gets +%d points for constructive comment\n", username, cfg.ConstructiveCommentPoint)
			}
		}
	}

	// Generate leaderboard
	leaderboard := karma.GenerateLeaderboard(reviewerKarma)

	// Write leaderboard to file with custom scoring display
	err = karma.WriteLeaderboardFileWithConfig(leaderboard, cfg.ReviewPoint, cfg.PositiveEmojiPoint, cfg.ConstructiveCommentPoint)
	if err != nil {
		fmt.Printf("❌ Error writing leaderboard file: %v\n", err)
		os.Exit(1)
	}
}

func runIncrementalUpdate(ctx context.Context, client *github.Client, owner, repo string, cfg config.Config) {
	fmt.Println("🔄 Running in incremental update mode...")

	// Initialize storage
	storage := storage.NewStorage(".karma-data.json")

	// Load existing karma data
	karmaData, err := storage.GetKarmaData()
	if err != nil {
		fmt.Printf("⚠️ Error loading karma data: %v\n", err)
		fmt.Println("🔄 Starting fresh...")
		// Create empty data using the storage instance
		karmaData, err = storage.Load()
		if err != nil {
			fmt.Printf("⚠️ Error creating empty karma data: %v\n", err)
			os.Exit(1)
		}
	}

	// Fetch all pull requests
	prs, err := githubapi.FetchAllPullRequests(ctx, client, owner, repo)
	if err != nil {
		fmt.Printf("❌ Error fetching pull requests: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📋 Found %d pull requests\n", len(prs))

	// Get processed PRs
	processedPRs, err := storage.GetProcessedPRs()
	if err != nil {
		fmt.Printf("⚠️ Error getting processed PRs: %v\n", err)
		processedPRs = make(map[int]bool)
	}

	// Process only new PRs
	newPRsCount := 0
	for _, pr := range prs {
		if processedPRs[pr.GetNumber()] {
			continue // Skip already processed PRs
		}

		newPRsCount++
		fmt.Printf("🆕 Processing new PR #%d: %s\n", pr.GetNumber(), pr.GetTitle())

		// Calculate karma for this PR
		prKarma := calculatePRKarma(ctx, client, owner, repo, pr.GetNumber(), cfg)

		// Update storage
		err = storage.UpdateKarma(pr.GetNumber(), prKarma)
		if err != nil {
			fmt.Printf("⚠️ Error updating karma for PR #%d: %v\n", pr.GetNumber(), err)
			continue
		}

		// Update in-memory data
		for username, points := range prKarma {
			karmaData.Reviewers[username] += points
		}
	}

	if newPRsCount == 0 {
		fmt.Println("✅ No new PRs to process")
	} else {
		fmt.Printf("✅ Processed %d new PRs\n", newPRsCount)
	}

	// Generate leaderboard from updated data
	leaderboard := karma.GenerateLeaderboard(karmaData.Reviewers)

	// Write leaderboard to file with custom scoring display
	err = karma.WriteLeaderboardFileWithConfig(leaderboard, cfg.ReviewPoint, cfg.PositiveEmojiPoint, cfg.ConstructiveCommentPoint)
	if err != nil {
		fmt.Printf("❌ Error writing leaderboard file: %v\n", err)
		os.Exit(1)
	}
}

func calculatePRKarma(ctx context.Context, client *github.Client, owner, repo string, prNumber int, cfg config.Config) map[string]int {
	reviewerKarma := make(map[string]int)

	// Get reviews for this PR
	reviews, err := githubapi.FetchPullRequestReviews(ctx, client, owner, repo, prNumber)
	if err != nil {
		fmt.Printf("⚠️ Error fetching reviews for PR #%d: %v\n", prNumber, err)
		return reviewerKarma
	}

	// Process reviews
	for _, review := range reviews {
		username := review.GetUser().GetLogin()
		if karma.IsBot(username) {
			continue
		}

		// Award points for giving a review
		reviewerKarma[username] += cfg.ReviewPoint

		// Check for positive emojis in review body
		if karma.HasPositiveEmoji(review.GetBody()) {
			reviewerKarma[username] += cfg.PositiveEmojiPoint
			fmt.Printf("  🎉 @%s gets +%d points for positive emoji\n", username, cfg.PositiveEmojiPoint)
		}

		// Check for constructive comments
		if karma.IsConstructiveComment(review.GetBody()) {
			reviewerKarma[username] += cfg.ConstructiveCommentPoint
			fmt.Printf("  💬 @%s gets +%d points for constructive comment\n", username, cfg.ConstructiveCommentPoint)
		}
	}

	// Get comments for this PR
	comments, err := githubapi.FetchPullRequestComments(ctx, client, owner, repo, prNumber)
	if err != nil {
		fmt.Printf("⚠️ Error fetching comments for PR #%d: %v\n", prNumber, err)
		return reviewerKarma
	}

	// Process comments
	for _, comment := range comments {
		username := comment.GetUser().GetLogin()
		if karma.IsBot(username) {
			continue
		}

		// Check for positive emojis in comment body
		if karma.HasPositiveEmoji(comment.GetBody()) {
			reviewerKarma[username] += cfg.PositiveEmojiPoint
			fmt.Printf("  🎉 @%s gets +%d points for positive emoji in comment\n", username, cfg.PositiveEmojiPoint)
		}

		// Check for constructive comments
		if karma.IsConstructiveComment(comment.GetBody()) {
			reviewerKarma[username] += cfg.ConstructiveCommentPoint
			fmt.Printf("  💬 @%s gets +%d points for constructive comment\n", username, cfg.ConstructiveCommentPoint)
		}
	}

	return reviewerKarma
}

func getUpdateModeString(incremental bool) string {
	if incremental {
		return "Incremental (only new PRs)"
	}
	return "Full Recreation (all PRs)"
}
