package repositories

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Gust4voSales/appstore-rss-reviews-app/server/internal/models"
)

// TestMain suppresses logs during all tests
func TestMain(m *testing.M) {
	// Suppress logs during testing
	log.SetOutput(io.Discard)

	// Run tests
	code := m.Run()

	// Restore log output
	log.SetOutput(os.Stderr)

	// Exit with the same code as the tests
	os.Exit(code)
}

// Helper function to create test reviews with different timestamps
func createTestReviews() models.AppStoreReviews {
	now := time.Now().UTC()
	return models.AppStoreReviews{
		{
			ID:        "review-1",
			Title:     "Great App",
			Content:   "This app is amazing!",
			Author:    "User1",
			Rating:    5,
			UpdatedAt: now.Add(-1 * time.Hour), // 1 hour ago
		},
		{
			ID:        "review-2",
			Title:     "Good App",
			Content:   "Pretty good app overall",
			Author:    "User2",
			Rating:    4,
			UpdatedAt: now.Add(-3 * time.Hour), // 3 hours ago
		},
		{
			ID:        "review-3",
			Title:     "Average App",
			Content:   "It's okay, could be better",
			Author:    "User3",
			Rating:    3,
			UpdatedAt: now.Add(-25 * time.Hour), // 25 hours ago (more than 24 hours)
		},
		{
			ID:        "review-4",
			Title:     "Poor App",
			Content:   "Not satisfied with this app",
			Author:    "User4",
			Rating:    2,
			UpdatedAt: now.Add(-48 * time.Hour), // 48 hours ago (more than 24 hours)
		},
	}
}

// Helper function to create a temporary test file with JSON data
func createTempFileWithReviews(t *testing.T, reviews models.AppStoreReviews) string {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test_reviews.json")

	data, err := json.MarshalIndent(reviews, "", "	")
	if err != nil {
		t.Fatalf("Failed to marshal test reviews: %v", err)
	}

	err = os.WriteFile(filePath, data, os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	return filePath
}

// Helper function to create a temporary file with invalid JSON
func createTempFileWithInvalidJSON(t *testing.T) string {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "invalid_reviews.json")

	invalidJSON := `{"invalid": json content without proper closing`
	err := os.WriteFile(filePath, []byte(invalidJSON), os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to write invalid JSON file: %v", err)
	}

	return filePath
}

// TestLoad_WithExistingFile verifies that Load correctly loads reviews from an existing file
func TestLoad_WithExistingFile(t *testing.T) {
	testReviews := createTestReviews()
	filePath := createTempFileWithReviews(t, testReviews)

	repo := Load(filePath)

	if repo == nil {
		t.Fatal("Expected repository to be created, got nil")
	}

	if repo.StorageFilePath != filePath {
		t.Errorf("Expected StorageFilePath to be %s, got %s", filePath, repo.StorageFilePath)
	}

	if len(repo.Reviews) != len(testReviews) {
		t.Errorf("Expected %d reviews to be loaded, got %d", len(testReviews), len(repo.Reviews))
	}

	// Verify the reviews are loaded correctly
	for i, expectedReview := range testReviews {
		if i >= len(repo.Reviews) {
			t.Errorf("Missing review at index %d", i)
			continue
		}

		actualReview := repo.Reviews[i]
		if actualReview.ID != expectedReview.ID {
			t.Errorf("Expected review ID %s, got %s", expectedReview.ID, actualReview.ID)
		}
		if actualReview.Title != expectedReview.Title {
			t.Errorf("Expected review title %s, got %s", expectedReview.Title, actualReview.Title)
		}
	}
}

// TestLoad_WithNonExistentFile verifies that Load handles non-existent files gracefully
func TestLoad_WithNonExistentFile(t *testing.T) {
	nonExistentPath := filepath.Join(t.TempDir(), "nonexistent.json")

	repo := Load(nonExistentPath)

	if repo == nil {
		t.Fatal("Expected repository to be created, got nil")
	}

	if len(repo.Reviews) != 0 {
		t.Errorf("Expected empty reviews list for non-existent file, got %d reviews", len(repo.Reviews))
	}
}

// TestLoad_WithInvalidJSON verifies that Load handles invalid JSON gracefully
func TestLoad_WithInvalidJSON(t *testing.T) {
	filePath := createTempFileWithInvalidJSON(t)

	repo := Load(filePath)

	if repo == nil {
		t.Fatal("Expected repository to be created, got nil")
	}

	if repo.StorageFilePath != filePath {
		t.Errorf("Expected StorageFilePath to be %s, got %s", filePath, repo.StorageFilePath)
	}

	// Should start with empty reviews when JSON is invalid
	if len(repo.Reviews) != 0 {
		t.Errorf("Expected empty reviews list for invalid JSON, got %d reviews", len(repo.Reviews))
	}
}

// TestListLatest_WithVariousTimeRanges verifies ListLatest with different hour ranges
func TestListLatest_WithVariousTimeRanges(t *testing.T) {
	testReviews := createTestReviews()
	filePath := createTempFileWithReviews(t, testReviews)
	repo := Load(filePath)

	// Test with 2 hours - should return reviews from 1 hour ago
	recentReviews := repo.ListLatest(2, ReviewFilter{})
	if len(recentReviews) != 1 {
		t.Errorf("Expected 1 review within 2 hours, got %d", len(recentReviews))
	}
	if len(recentReviews) > 0 && recentReviews[0].ID != "review-1" {
		t.Errorf("Expected review-1 to be the most recent, got %s", recentReviews[0].ID)
	}

	// Test with 5 hours - should return reviews from 1 and 3 hours ago
	recentReviews = repo.ListLatest(5, ReviewFilter{})
	if len(recentReviews) != 2 {
		t.Errorf("Expected 2 reviews within 5 hours, got %d", len(recentReviews))
	}

	// Test with 24 hours - should return reviews from 1 and 3 hours ago (not 25 hours ago)
	recentReviews = repo.ListLatest(24, ReviewFilter{})
	if len(recentReviews) != 2 {
		t.Errorf("Expected 2 reviews within 24 hours, got %d", len(recentReviews))
	}

	// Test with 50 hours - should return all reviews
	recentReviews = repo.ListLatest(50, ReviewFilter{})
	if len(recentReviews) != 4 {
		t.Errorf("Expected 4 reviews within 50 hours, got %d", len(recentReviews))
	}

	// Test with 0 hours - should return no reviews
	recentReviews = repo.ListLatest(0, ReviewFilter{})
	if len(recentReviews) != 0 {
		t.Errorf("Expected 0 reviews within 0 hours, got %d", len(recentReviews))
	}
}

// TestListLatest_WithEmptyRepository verifies ListLatest with empty repository
func TestListLatest_WithEmptyRepository(t *testing.T) {
	repo := Load("")

	recentReviews := repo.ListLatest(24, ReviewFilter{})

	if len(recentReviews) != 0 {
		t.Errorf("Expected 0 reviews from empty repository, got %d", len(recentReviews))
	}
}

// TestGetLatestReview_WithReviews verifies GetLatestReview returns the most recent review
func TestGetLatestReview_WithReviews(t *testing.T) {
	testReviews := createTestReviews()
	filePath := createTempFileWithReviews(t, testReviews)
	repo := Load(filePath)

	latestReview := repo.GetLatestReview()

	if latestReview == nil {
		t.Fatal("Expected latest review to be returned, got nil")
	}

	// Since reviews are sorted by UpdatedAt descending, the first one should be the latest
	expectedID := testReviews[0].ID
	if latestReview.ID != expectedID {
		t.Errorf("Expected latest review ID to be %s, got %s", expectedID, latestReview.ID)
	}
}

// TestGetLatestReview_WithEmptyRepository verifies GetLatestReview returns nil for empty repository
func TestGetLatestReview_WithEmptyRepository(t *testing.T) {
	repo := Load("")

	latestReview := repo.GetLatestReview()

	if latestReview != nil {
		t.Errorf("Expected nil for empty repository, got review with ID %s", latestReview.ID)
	}
}

// TestAddBatch_WithNewReviews verifies AddBatch adds new reviews and returns correct count
func TestAddBatch_WithNewReviews(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test_add_batch.json")
	repo := Load(filePath)

	newReviews := models.AppStoreReviews{
		{
			ID:        "new-review-1",
			Title:     "New Review 1",
			Content:   "Content 1",
			Author:    "Author1",
			Rating:    5,
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:        "new-review-2",
			Title:     "New Review 2",
			Content:   "Content 2",
			Author:    "Author2",
			Rating:    4,
			UpdatedAt: time.Now().UTC().Add(-1 * time.Hour),
		},
	}

	addedCount, err := repo.AddBatch(newReviews)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if addedCount != 2 {
		t.Errorf("Expected 2 reviews to be added, got %d", addedCount)
	}

	if len(repo.Reviews) != 2 {
		t.Errorf("Expected repository to have 2 reviews, got %d", len(repo.Reviews))
	}

	// Verify reviews are sorted by UpdatedAt descending
	if repo.Reviews[0].ID != "new-review-1" {
		t.Errorf("Expected first review to be new-review-1 (most recent), got %s", repo.Reviews[0].ID)
	}
}

// TestAddBatch_WithDuplicateReviews verifies AddBatch handles duplicate reviews correctly
func TestAddBatch_WithDuplicateReviews(t *testing.T) {
	testReviews := createTestReviews()
	filePath := createTempFileWithReviews(t, testReviews)
	repo := Load(filePath)

	initialCount := len(repo.Reviews)

	// Try to add reviews with some duplicates
	mixedReviews := models.AppStoreReviews{
		{
			ID:        "review-1", // Duplicate
			Title:     "Updated Great App",
			Content:   "This app is still amazing!",
			Author:    "User1",
			Rating:    5,
			UpdatedAt: time.Now().UTC(),
		},
		{
			ID:        "new-unique-review",
			Title:     "Brand New Review",
			Content:   "This is a completely new review",
			Author:    "NewUser",
			Rating:    4,
			UpdatedAt: time.Now().UTC().Add(-30 * time.Minute),
		},
	}

	addedCount, err := repo.AddBatch(mixedReviews)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if addedCount != 1 {
		t.Errorf("Expected 1 new review to be added (excluding duplicate), got %d", addedCount)
	}

	expectedTotalCount := initialCount + 1
	if len(repo.Reviews) != expectedTotalCount {
		t.Errorf("Expected repository to have %d reviews, got %d", expectedTotalCount, len(repo.Reviews))
	}
}

// TestAddBatch_SortingAfterAdd verifies that reviews are properly sorted after adding batch
func TestAddBatch_SortingAfterAdd(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test_sorting.json")
	repo := Load(filePath)

	now := time.Now().UTC()
	unsortedReviews := models.AppStoreReviews{
		{
			ID:        "old-review",
			Title:     "Old Review",
			Content:   "This is an old review",
			Author:    "OldUser",
			Rating:    3,
			UpdatedAt: now.Add(-2 * time.Hour), // 2 hours ago
		},
		{
			ID:        "newest-review",
			Title:     "Newest Review",
			Content:   "This is the newest review",
			Author:    "NewestUser",
			Rating:    5,
			UpdatedAt: now, // Most recent
		},
		{
			ID:        "middle-review",
			Title:     "Middle Review",
			Content:   "This is a middle review",
			Author:    "MiddleUser",
			Rating:    4,
			UpdatedAt: now.Add(-1 * time.Hour), // 1 hour ago
		},
	}

	_, err := repo.AddBatch(unsortedReviews)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify reviews are sorted by UpdatedAt descending
	expectedOrder := []string{"newest-review", "middle-review", "old-review"}
	for i, expectedID := range expectedOrder {
		if i >= len(repo.Reviews) {
			t.Errorf("Missing review at index %d", i)
			continue
		}
		if repo.Reviews[i].ID != expectedID {
			t.Errorf("Expected review at index %d to be %s, got %s", i, expectedID, repo.Reviews[i].ID)
		}
	}
}

// TestHasReviewWithID verifies the hasReviewWithID helper method
func TestHasReviewWithID(t *testing.T) {
	testReviews := createTestReviews()
	repo := &AppReviewsRepository{
		Reviews: testReviews,
	}

	// Test existing ID
	if !repo.hasReviewWithID("review-1") {
		t.Error("Expected hasReviewWithID to return true for existing review-1")
	}

	// Test non-existing ID
	if repo.hasReviewWithID("non-existent-id") {
		t.Error("Expected hasReviewWithID to return false for non-existent ID")
	}

	// Test with empty repository
	emptyRepo := &AppReviewsRepository{
		Reviews: models.AppStoreReviews{},
	}
	if emptyRepo.hasReviewWithID("any-id") {
		t.Error("Expected hasReviewWithID to return false for empty repository")
	}
}

// TestAddBatch_FileOperationError verifies AddBatch handles file operation errors
func TestAddBatch_FileOperationError(t *testing.T) {
	// Create a directory where we expect a file - this should cause a write error
	tmpDir := t.TempDir()
	dirPath := filepath.Join(tmpDir, "should_be_file")
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create directory for test: %v", err)
	}

	repo := &AppReviewsRepository{
		Reviews:         models.AppStoreReviews{},
		StorageFilePath: dirPath, // directory, not a file
	}

	newReviews := models.AppStoreReviews{
		{
			ID:        "test-review",
			Title:     "Test Review",
			Content:   "Test content",
			Author:    "Test Author",
			Rating:    5,
			UpdatedAt: time.Now().UTC(),
		},
	}

	addedCount, err := repo.AddBatch(newReviews)
	if err == nil {
		t.Error("Expected error when trying to write to directory path, got nil")
	}

	if addedCount != 0 {
		t.Errorf("Expected 0 reviews to be added when file operation fails, got %d", addedCount)
	}
}

// TestListLatest_WithRatingFilter verifies ListLatest correctly filters reviews by rating
func TestListLatest_WithRatingFilter(t *testing.T) {
	testReviews := createTestReviews()
	filePath := createTempFileWithReviews(t, testReviews)
	repo := Load(filePath)

	// Test filtering by 2-star rating - should return only review-4
	rating2 := 2
	recentReviews := repo.ListLatest(50, ReviewFilter{Rating: &rating2})
	if len(recentReviews) != 1 {
		t.Errorf("Expected 1 review with 2-star rating, got %d", len(recentReviews))
	}
	if len(recentReviews) > 0 && recentReviews[0].ID != "review-4" {
		t.Errorf("Expected review-4 (2-star review), got %s", recentReviews[0].ID)
	}

	// Test filtering by 1-star rating - should return no reviews (none exist)
	rating1 := 1
	recentReviews = repo.ListLatest(50, ReviewFilter{Rating: &rating1})
	if len(recentReviews) != 0 {
		t.Errorf("Expected 0 reviews with 1-star rating, got %d", len(recentReviews))
	}

	// Test filtering by rating combined with time constraints
	// Only reviews within 5 hours (review-1 and review-2) with 4-star rating (review-2)
	rating4 := 4
	recentReviews = repo.ListLatest(5, ReviewFilter{Rating: &rating4})
	if len(recentReviews) != 1 {
		t.Errorf("Expected 1 review with 4-star rating within 5 hours, got %d", len(recentReviews))
	}
	if len(recentReviews) > 0 && recentReviews[0].ID != "review-2" {
		t.Errorf("Expected review-2 (4-star review within 5 hours), got %s", recentReviews[0].ID)
	}

	// Test filtering by rating with time constraints that exclude the matching review
	// Only reviews within 2 hours (review-1) with 4-star rating - should return nothing
	recentReviews = repo.ListLatest(2, ReviewFilter{Rating: &rating4})
	if len(recentReviews) != 0 {
		t.Errorf("Expected 0 reviews with 4-star rating within 2 hours, got %d", len(recentReviews))
	}
}
