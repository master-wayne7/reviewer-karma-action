package karma

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

// Reviewer represents a user with their karma points
type Reviewer struct {
	Username string `json:"username"`
	Points   int    `json:"points"`
}

// Leaderboard represents the karma leaderboard
type Leaderboard struct {
	Reviewers []Reviewer `json:"reviewers"`
}

// Positive emojis that award bonus points
var positiveEmojis = map[string]bool{
	"👍": true, "🔥": true, "😄": true, "🎉": true, "🚀": true,
	"💯": true, "✅": true, "⭐": true, "❤️": true, "👏": true,
}

// IsBot checks if a username belongs to a bot
func IsBot(username string) bool {
	botSuffixes := []string{"[bot]", "-bot", "bot-", "github-actions[bot]", "dependabot[bot]"}
	usernameLower := strings.ToLower(username)

	for _, suffix := range botSuffixes {
		if strings.Contains(usernameLower, strings.ToLower(suffix)) {
			return true
		}
	}

	return false
}

// HasPositiveEmoji checks if text contains positive emojis
func HasPositiveEmoji(text string) bool {
	if text == "" {
		return false
	}

	for emoji := range positiveEmojis {
		if strings.Contains(text, emoji) {
			return true
		}
	}

	return false
}

// IsConstructiveComment checks if a comment is constructive
func IsConstructiveComment(text string) bool {
	if text == "" {
		return false
	}

	// Remove common non-constructive phrases
	text = strings.ToLower(text)
	nonConstructive := []string{
		"lgtm", "looks good", "good", "nice", "👍", "🔥", "😄", "🎉", "🚀", "💯", "✅", "⭐", "❤️", "👏",
	}

	for _, phrase := range nonConstructive {
		// Only remove exact matches or standalone words
		text = strings.ReplaceAll(text, " "+phrase+" ", " ")
		text = strings.ReplaceAll(text, phrase+" ", " ")
		text = strings.ReplaceAll(text, " "+phrase, " ")
		text = strings.ReplaceAll(text, phrase, "")
	}

	// Count remaining words
	words := strings.Fields(text)
	return len(words) > 10
}

// GenerateLeaderboard creates a leaderboard from reviewer karma
func GenerateLeaderboard(reviewerKarma map[string]int) Leaderboard {
	var reviewers []Reviewer

	for username, points := range reviewerKarma {
		reviewers = append(reviewers, Reviewer{
			Username: username,
			Points:   points,
		})
	}

	// Sort by points (descending)
	sort.Slice(reviewers, func(i, j int) bool {
		return reviewers[i].Points > reviewers[j].Points
	})

	return Leaderboard{Reviewers: reviewers}
}

// WriteLeaderboardFile writes the leaderboard to REVIEWERS.md
func WriteLeaderboardFile(leaderboard Leaderboard) error {
	content := generateLeaderboardMarkdown(leaderboard)

	// Write to REVIEWERS.md
	err := os.WriteFile("REVIEWERS.md", []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}

// WriteLeaderboardFileWithConfig writes the leaderboard to REVIEWERS.md with custom scoring display
func WriteLeaderboardFileWithConfig(leaderboard Leaderboard, reviewPoint, emojiPoint, commentPoint int) error {
	content := generateLeaderboardMarkdownWithConfig(leaderboard, reviewPoint, emojiPoint, commentPoint)

	// Write to REVIEWERS.md
	err := os.WriteFile("REVIEWERS.md", []byte(content), 0644)
	if err != nil {
		return err
	}

	return nil
}

// generateLeaderboardMarkdown generates markdown content for the leaderboard
func generateLeaderboardMarkdown(leaderboard Leaderboard) string {
	return generateLeaderboardMarkdownWithConfig(leaderboard, 1, 2, 1) // Default values
}

// generateLeaderboardMarkdownWithConfig generates markdown content with custom scoring display
func generateLeaderboardMarkdownWithConfig(leaderboard Leaderboard, reviewPoint, emojiPoint, commentPoint int) string {
	var sb strings.Builder

	sb.WriteString("# Reviewer Karma Leaderboard\n\n")
	sb.WriteString("This leaderboard tracks reviewer engagement and contributions to the repository.\n\n")
	sb.WriteString("## Scoring System\n\n")
	sb.WriteString(fmt.Sprintf("- ✅ Giving a code review: +%d point(s)\n", reviewPoint))
	sb.WriteString(fmt.Sprintf("- ✅ Review includes a positive emoji (👍, 🔥, 😄, etc.): +%d point(s)\n", emojiPoint))
	sb.WriteString(fmt.Sprintf("- ✅ Review comment contains a constructive message (>10 words): +%d point(s)\n\n", commentPoint))
	sb.WriteString("## Current Rankings\n\n")
	sb.WriteString("| Rank | Reviewer | Points |\n")
	sb.WriteString("|------|----------|--------|\n")

	medals := []string{"🥇", "🥈", "🥉"}

	for i, reviewer := range leaderboard.Reviewers {
		rank := i + 1
		medal := ""

		if rank <= 3 {
			medal = medals[rank-1] + " "
		}

		sb.WriteString(fmt.Sprintf("| %d | %s@%s | %d |\n", rank, medal, reviewer.Username, reviewer.Points))
	}

	sb.WriteString("\n---\n")
	sb.WriteString(fmt.Sprintf("*Last updated: %s*\n", time.Now().Format("2006-01-02 15:04:05 UTC")))

	return sb.String()
}
