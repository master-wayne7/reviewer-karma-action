package storage

import (
	"os"
	"testing"
	"time"
)

func TestStorage_LoadSave(t *testing.T) {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "karma_test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	storage := NewStorage(tmpFile.Name())

	// Test loading empty data
	data, err := storage.Load()
	if err != nil {
		t.Fatalf("Failed to load empty data: %v", err)
	}

	if data.Reviewers == nil {
		t.Error("Reviewers map should be initialized")
	}
	if data.ProcessedPRs == nil {
		t.Error("ProcessedPRs map should be initialized")
	}

	// Test saving and loading data
	testData := &KarmaData{
		Reviewers: map[string]int{
			"alice": 10,
			"bob":   5,
		},
		ProcessedPRs: map[int]time.Time{
			1: time.Now(),
			2: time.Now(),
		},
	}

	err = storage.Save(testData)
	if err != nil {
		t.Fatalf("Failed to save data: %v", err)
	}

	loadedData, err := storage.Load()
	if err != nil {
		t.Fatalf("Failed to load data: %v", err)
	}

	if len(loadedData.Reviewers) != 2 {
		t.Errorf("Expected 2 reviewers, got %d", len(loadedData.Reviewers))
	}

	if loadedData.Reviewers["alice"] != 10 {
		t.Errorf("Expected alice to have 10 points, got %d", loadedData.Reviewers["alice"])
	}

	if len(loadedData.ProcessedPRs) != 2 {
		t.Errorf("Expected 2 processed PRs, got %d", len(loadedData.ProcessedPRs))
	}
}

func TestStorage_UpdateKarma(t *testing.T) {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "karma_test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	storage := NewStorage(tmpFile.Name())

	// Test updating karma
	reviewerKarma := map[string]int{
		"alice": 5,
		"bob":   3,
	}

	err = storage.UpdateKarma(1, reviewerKarma)
	if err != nil {
		t.Fatalf("Failed to update karma: %v", err)
	}

	// Test updating karma again
	reviewerKarma2 := map[string]int{
		"alice": 2,
		"carol": 4,
	}

	err = storage.UpdateKarma(2, reviewerKarma2)
	if err != nil {
		t.Fatalf("Failed to update karma again: %v", err)
	}

	// Verify final state
	data, err := storage.Load()
	if err != nil {
		t.Fatalf("Failed to load data: %v", err)
	}

	expectedAlice := 7 // 5 + 2
	expectedBob := 3   // 3
	expectedCarol := 4 // 4

	if data.Reviewers["alice"] != expectedAlice {
		t.Errorf("Expected alice to have %d points, got %d", expectedAlice, data.Reviewers["alice"])
	}

	if data.Reviewers["bob"] != expectedBob {
		t.Errorf("Expected bob to have %d points, got %d", expectedBob, data.Reviewers["bob"])
	}

	if data.Reviewers["carol"] != expectedCarol {
		t.Errorf("Expected carol to have %d points, got %d", expectedCarol, data.Reviewers["carol"])
	}

	// Verify processed PRs
	processedPRs, err := storage.GetProcessedPRs()
	if err != nil {
		t.Fatalf("Failed to get processed PRs: %v", err)
	}

	if !processedPRs[1] {
		t.Error("PR 1 should be marked as processed")
	}

	if !processedPRs[2] {
		t.Error("PR 2 should be marked as processed")
	}
}

func TestStorage_Clear(t *testing.T) {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "karma_test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	storage := NewStorage(tmpFile.Name())

	// Add some data
	reviewerKarma := map[string]int{"alice": 5}
	err = storage.UpdateKarma(1, reviewerKarma)
	if err != nil {
		t.Fatalf("Failed to update karma: %v", err)
	}

	// Clear data
	err = storage.Clear()
	if err != nil {
		t.Fatalf("Failed to clear data: %v", err)
	}

	// Verify data is cleared
	data, err := storage.Load()
	if err != nil {
		t.Fatalf("Failed to load data: %v", err)
	}

	if len(data.Reviewers) != 0 {
		t.Errorf("Expected 0 reviewers after clear, got %d", len(data.Reviewers))
	}

	if len(data.ProcessedPRs) != 0 {
		t.Errorf("Expected 0 processed PRs after clear, got %d", len(data.ProcessedPRs))
	}
}
