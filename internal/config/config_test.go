package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Test default config
	config := Load()

	if config.ReviewPoint != 1 {
		t.Errorf("Expected ReviewPoint to be 1, got %d", config.ReviewPoint)
	}

	if config.PositiveEmojiPoint != 2 {
		t.Errorf("Expected PositiveEmojiPoint to be 2, got %d", config.PositiveEmojiPoint)
	}

	if config.ConstructiveCommentPoint != 1 {
		t.Errorf("Expected ConstructiveCommentPoint to be 1, got %d", config.ConstructiveCommentPoint)
	}
}

func TestLoadConfigWithEnvironment(t *testing.T) {
	// Set environment variables
	os.Setenv("REVIEW_POINT", "3")
	os.Setenv("POSITIVE_EMOJI_POINT", "5")
	os.Setenv("CONSTRUCTIVE_COMMENT_POINT", "2")

	// Test config with environment variables
	config := Load()

	if config.ReviewPoint != 3 {
		t.Errorf("Expected ReviewPoint to be 3, got %d", config.ReviewPoint)
	}

	if config.PositiveEmojiPoint != 5 {
		t.Errorf("Expected PositiveEmojiPoint to be 5, got %d", config.PositiveEmojiPoint)
	}

	if config.ConstructiveCommentPoint != 2 {
		t.Errorf("Expected ConstructiveCommentPoint to be 2, got %d", config.ConstructiveCommentPoint)
	}

	// Clean up
	os.Unsetenv("REVIEW_POINT")
	os.Unsetenv("POSITIVE_EMOJI_POINT")
	os.Unsetenv("CONSTRUCTIVE_COMMENT_POINT")
}
