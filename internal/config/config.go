package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds the karma point configuration
type Config struct {
	ReviewPoint              int
	PositiveEmojiPoint       int
	ConstructiveCommentPoint int
	IncrementalUpdate        bool
}

// Default configuration
var defaultConfig = Config{
	ReviewPoint:              1,
	PositiveEmojiPoint:       2,
	ConstructiveCommentPoint: 1,
	IncrementalUpdate:        false, // Default to full recreation
}

// Load loads configuration from environment variables
func Load() Config {
	config := defaultConfig

	if val := os.Getenv("REVIEW_POINT"); val != "" {
		if points, err := strconv.Atoi(val); err == nil {
			config.ReviewPoint = points
		}
	}

	if val := os.Getenv("POSITIVE_EMOJI_POINT"); val != "" {
		if points, err := strconv.Atoi(val); err == nil {
			config.PositiveEmojiPoint = points
		}
	}

	if val := os.Getenv("CONSTRUCTIVE_COMMENT_POINT"); val != "" {
		if points, err := strconv.Atoi(val); err == nil {
			config.ConstructiveCommentPoint = points
		}
	}

	if val := os.Getenv("INCREMENTAL_UPDATE"); val != "" {
		config.IncrementalUpdate = strings.ToLower(val) == "true"
	}

	return config
}
