package karma

import (
	"strings"
	"testing"
)

func TestIsBot(t *testing.T) {
	tests := []struct {
		username string
		expected bool
	}{
		{"alice", false},
		{"bob", false},
		{"github-actions[bot]", true},
		{"dependabot[bot]", true},
		{"test-bot", true},
		{"bot-user", true},
		{"user-bot", true},
		{"normaluser", false},
	}

	for _, test := range tests {
		result := IsBot(test.username)
		if result != test.expected {
			t.Errorf("IsBot(%s) = %v, expected %v", test.username, result, test.expected)
		}
	}
}

func TestHasPositiveEmoji(t *testing.T) {
	tests := []struct {
		text     string
		expected bool
	}{
		{"", false},
		{"This is a normal comment", false},
		{"Great work! 👍", true},
		{"Amazing! 🔥", true},
		{"Nice job 😄", true},
		{"Good work 🎉", true},
		{"Excellent 🚀", true},
		{"Perfect 💯", true},
		{"Looks good ✅", true},
		{"Awesome ⭐", true},
		{"Love it ❤️", true},
		{"Well done 👏", true},
		{"This is great but no emoji", false},
	}

	for _, test := range tests {
		result := HasPositiveEmoji(test.text)
		if result != test.expected {
			t.Errorf("HasPositiveEmoji(%q) = %v, expected %v", test.text, result, test.expected)
		}
	}
}

func TestIsConstructiveComment(t *testing.T) {
	tests := []struct {
		text     string
		expected bool
	}{
		{"", false},
		{"LGTM", false},
		{"Looks good", false},
		{"Good 👍", false},
		{"Nice 🔥", false},
		{"This is a very detailed comment that provides constructive feedback about the code changes and suggests improvements for better maintainability", true},
		{"I think we should refactor this function to improve readability and add better error handling", true},
		{"The implementation looks good but we should consider adding more test cases", false}, // Exactly 10 words after filtering
		{"LGTM but we should add more documentation", false},                                   // Contains LGTM
		{"Great work! 👍 This is excellent", false},                                             // Contains emoji
	}

	for _, test := range tests {
		result := IsConstructiveComment(test.text)
		if result != test.expected {
			t.Errorf("IsConstructiveComment(%q) = %v, expected %v", test.text, result, test.expected)
		}
	}
}

func TestGenerateLeaderboard(t *testing.T) {
	reviewerKarma := map[string]int{
		"alice": 18,
		"bob":   12,
		"carol": 10,
		"dave":  8,
		"eve":   5,
	}

	leaderboard := GenerateLeaderboard(reviewerKarma)

	if len(leaderboard.Reviewers) != 5 {
		t.Errorf("Expected 5 reviewers, got %d", len(leaderboard.Reviewers))
	}

	// Check sorting (should be descending by points)
	for i := 0; i < len(leaderboard.Reviewers)-1; i++ {
		if leaderboard.Reviewers[i].Points < leaderboard.Reviewers[i+1].Points {
			t.Errorf("Leaderboard not sorted correctly: %d < %d",
				leaderboard.Reviewers[i].Points, leaderboard.Reviewers[i+1].Points)
		}
	}

	// Check first place
	if leaderboard.Reviewers[0].Username != "alice" || leaderboard.Reviewers[0].Points != 18 {
		t.Errorf("Expected alice with 18 points, got %s with %d points",
			leaderboard.Reviewers[0].Username, leaderboard.Reviewers[0].Points)
	}
}

func TestDebugConstructiveComment(t *testing.T) {
	text := "The implementation looks good but we should consider adding more test cases"

	// Debug the processing
	textLower := strings.ToLower(text)
	t.Logf("Original: %q", text)
	t.Logf("Lowercase: %q", textLower)

	nonConstructive := []string{
		"lgtm", "looks good", "good", "nice", "👍", "🔥", "😄", "🎉", "🚀", "💯", "✅", "⭐", "❤️", "👏",
	}

	for _, phrase := range nonConstructive {
		before := textLower
		textLower = strings.ReplaceAll(textLower, " "+phrase+" ", " ")
		textLower = strings.ReplaceAll(textLower, phrase+" ", " ")
		textLower = strings.ReplaceAll(textLower, " "+phrase, " ")
		textLower = strings.ReplaceAll(textLower, phrase, "")
		if before != textLower {
			t.Logf("After removing '%s': %q", phrase, textLower)
		}
	}

	words := strings.Fields(textLower)
	t.Logf("Final words: %v (count: %d)", words, len(words))
}
