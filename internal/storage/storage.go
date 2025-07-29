package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// KarmaData represents the stored karma data
type KarmaData struct {
	Reviewers    map[string]int    `json:"reviewers"`
	LastUpdated  time.Time         `json:"last_updated"`
	ProcessedPRs map[int]time.Time `json:"processed_prs"` // PR number -> last processed time
}

// Storage handles persistence of karma data
type Storage struct {
	filePath string
}

// NewStorage creates a new storage instance
func NewStorage(filePath string) *Storage {
	return &Storage{
		filePath: filePath,
	}
}

// Load loads karma data from file
func (s *Storage) Load() (*KarmaData, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, return empty data
			return &KarmaData{
				Reviewers:    make(map[string]int),
				LastUpdated:  time.Time{},
				ProcessedPRs: make(map[int]time.Time),
			}, nil
		}
		return nil, fmt.Errorf("failed to read karma data: %w", err)
	}

	// Handle empty file
	if len(data) == 0 {
		return &KarmaData{
			Reviewers:    make(map[string]int),
			LastUpdated:  time.Time{},
			ProcessedPRs: make(map[int]time.Time),
		}, nil
	}

	var karmaData KarmaData
	if err := json.Unmarshal(data, &karmaData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal karma data: %w", err)
	}

	// Initialize maps if they're nil
	if karmaData.Reviewers == nil {
		karmaData.Reviewers = make(map[string]int)
	}
	if karmaData.ProcessedPRs == nil {
		karmaData.ProcessedPRs = make(map[int]time.Time)
	}

	return &karmaData, nil
}

// Save saves karma data to file
func (s *Storage) Save(data *KarmaData) error {
	data.LastUpdated = time.Now()

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal karma data: %w", err)
	}

	if err := os.WriteFile(s.filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write karma data: %w", err)
	}

	return nil
}

// UpdateKarma updates karma data with new reviewer points
func (s *Storage) UpdateKarma(prNumber int, reviewerKarma map[string]int) error {
	data, err := s.Load()
	if err != nil {
		return err
	}

	// Update reviewer karma
	for username, points := range reviewerKarma {
		data.Reviewers[username] += points
	}

	// Mark PR as processed
	data.ProcessedPRs[prNumber] = time.Now()

	return s.Save(data)
}

// GetProcessedPRs returns a map of processed PR numbers
func (s *Storage) GetProcessedPRs() (map[int]bool, error) {
	data, err := s.Load()
	if err != nil {
		return nil, err
	}

	processed := make(map[int]bool)
	for prNumber := range data.ProcessedPRs {
		processed[prNumber] = true
	}

	return processed, nil
}

// GetKarmaData returns the current karma data
func (s *Storage) GetKarmaData() (*KarmaData, error) {
	return s.Load()
}

// Clear clears all stored data
func (s *Storage) Clear() error {
	data := &KarmaData{
		Reviewers:    make(map[string]int),
		LastUpdated:  time.Now(),
		ProcessedPRs: make(map[int]time.Time),
	}
	return s.Save(data)
}

// NewEmptyKarmaData creates a new empty KarmaData instance
func NewEmptyKarmaData() *KarmaData {
	return &KarmaData{
		Reviewers:    make(map[string]int),
		LastUpdated:  time.Now(),
		ProcessedPRs: make(map[int]time.Time),
	}
}

// CreateEmptyKarmaData creates a new empty KarmaData instance
func CreateEmptyKarmaData() *KarmaData {
	return &KarmaData{
		Reviewers:    make(map[string]int),
		LastUpdated:  time.Now(),
		ProcessedPRs: make(map[int]time.Time),
	}
}
