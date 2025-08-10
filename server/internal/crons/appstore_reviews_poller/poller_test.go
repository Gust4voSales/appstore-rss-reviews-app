package appstore_reviews_poller

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Gust4voSales/appstore-rss-reviews-app/server/config"
	"github.com/Gust4voSales/appstore-rss-reviews-app/server/internal/models"
)

// MockApp is a mock implementation of the App service for testing
type MockApp struct {
	addReviewsFuncCalled int
	mockedLatestReview   *models.AppStoreReview
}

func (a *MockApp) ListLatestReviews(hours int, rating *int) []models.AppStoreReview {
	return nil
}
func (a *MockApp) GetLatestReview() *models.AppStoreReview {
	return a.mockedLatestReview
}
func (a *MockApp) AddReviews(reviews []models.AppStoreReview) (int, error) {
	a.addReviewsFuncCalled++
	return len(reviews), nil
}
func (a *MockApp) GetAppID() string {
	return ""
}

// MockFetcher is a mock implementation of the Fetcher for testing
type MockFetcher struct {
	mockedReviews          []models.AppStoreReview
	mockedError            error
	capturedAppID          string
	capturedLatestReviewId *string
	fetchReviewsCalled     int
}

func (f *MockFetcher) fetchReviews(ctx context.Context, appID string, latestReviewId *string) ([]models.AppStoreReview, error) {
	f.fetchReviewsCalled++
	f.capturedAppID = appID
	f.capturedLatestReviewId = latestReviewId

	if f.mockedError != nil {
		return nil, f.mockedError
	}
	return f.mockedReviews, nil
}

// createTestPoller creates a poller with mocked dependencies for testing
func createTestPoller(mockApp *MockApp, mockFetcher *MockFetcher) *AppStoreReviewsPoller {
	cfg := &config.Config{
		AppID:           "test-app-id",
		PollingInterval: 100 * time.Millisecond, // Short interval for testing
	}

	poller := &AppStoreReviewsPoller{
		cfg:        cfg,
		appService: mockApp,
		fetcher:    mockFetcher,
	}

	return poller
}

// TestRun_StopsGracefullyWhenContextCancelled verifies that the poller stops gracefully when context is cancelled
func TestRun_StopsGracefullyWhenContextCancelled(t *testing.T) {
	mockApp := &MockApp{}
	mockFetcher := &MockFetcher{}
	poller := createTestPoller(mockApp, mockFetcher)

	// Create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Track when the Run method completes
	done := make(chan bool)
	go func() {
		poller.Run(ctx)
		done <- true
	}()

	// Let it run for a short time to ensure it starts
	time.Sleep(50 * time.Millisecond)

	// Cancel the context
	cancel()

	// Wait for the poller to stop, with a timeout
	select {
	case <-done:
		// Success - poller stopped
	case <-time.After(1 * time.Second):
		t.Fatal("Poller did not stop within timeout after context cancellation")
	}
}

// TestRun_ExecutesImmediatelyBeforeTickerStarts verifies that the first run happens immediately before ticker starts
func TestRun_ExecutesImmediatelyBeforeTickerStarts(t *testing.T) {
	mockApp := &MockApp{}
	mockFetcher := &MockFetcher{
		mockedReviews: []models.AppStoreReview{
			{
				ID:        "test-review-1",
				Title:     "Test Review 1",
				Content:   "Test content 1",
				Author:    "Test Author 1",
				Rating:    5,
				UpdatedAt: time.Now(),
			},
		},
	}
	poller := createTestPoller(mockApp, mockFetcher)

	// Use a very long polling interval to ensure we're testing immediate execution
	poller.cfg.PollingInterval = 10 * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the poller in a goroutine
	go poller.Run(ctx)

	// Wait a short time for the immediate run to complete
	time.Sleep(50 * time.Millisecond)

	// Cancel to stop the poller
	cancel()

	// Give it time to stop gracefully
	time.Sleep(10 * time.Millisecond)

	// Verify that processLatestReviews was called immediately
	addCalls := mockApp.addReviewsFuncCalled
	if addCalls != 1 {
		t.Errorf("Expected exactly 1 call to AddReviews (immediate run), got %d", addCalls)
	}
}

// TestRun_ExecutesPeriodicProcessingWithTicker verifies that periodic processing works with ticker
func TestRun_ExecutesPeriodicProcessingWithTicker(t *testing.T) {
	mockApp := &MockApp{}
	mockFetcher := &MockFetcher{
		mockedReviews: []models.AppStoreReview{
			{
				ID:        "test-review-1",
				Title:     "Test Review 1",
				Content:   "Test content 1",
				Author:    "Test Author 1",
				Rating:    5,
				UpdatedAt: time.Now(),
			},
		},
	}
	poller := createTestPoller(mockApp, mockFetcher)

	// Use a short polling interval for testing
	poller.cfg.PollingInterval = 50 * time.Millisecond

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the poller in a goroutine
	go poller.Run(ctx)

	// Wait for multiple polling cycles
	time.Sleep(200 * time.Millisecond)

	// Cancel to stop the poller
	cancel()

	// Give it time to stop gracefully
	time.Sleep(20 * time.Millisecond)

	// Verify that processLatestReviews was called multiple times
	// Should be at least: 1 (immediate) + 3-4 (periodic calls in 200ms with 50ms interval)
	addCalls := mockApp.addReviewsFuncCalled
	if addCalls < 3 {
		t.Errorf("Expected at least 3 calls to AddReviews (1 immediate + periodic), got %d", addCalls)
	}

	// Verify we didn't call too many times (should not exceed reasonable bounds)
	// With 200ms runtime and 50ms interval, we expect at most 1 + 4 = 5 calls
	if addCalls > 6 {
		t.Errorf("Expected at most 6 calls to AddReviews, got %d (possible timing issue)", addCalls)
	}
}

// TestProcessLatestReviews_CallsAddReviewsWhenReviewsReturned verifies that AddReviews is called when fetcher returns reviews
func TestProcessLatestReviews_CallsAddReviewsWhenReviewsReturned(t *testing.T) {
	mockApp := &MockApp{}
	mockFetcher := &MockFetcher{
		mockedReviews: []models.AppStoreReview{
			{
				ID:        "test-review-1",
				Title:     "Test Review 1",
				Content:   "Test content 1",
				Author:    "Test Author 1",
				Rating:    5,
				UpdatedAt: time.Now(),
			},
			{
				ID:        "test-review-2",
				Title:     "Test Review 2",
				Content:   "Test content 2",
				Author:    "Test Author 2",
				Rating:    4,
				UpdatedAt: time.Now(),
			},
		},
	}
	poller := createTestPoller(mockApp, mockFetcher)

	ctx := context.Background()

	// Call processLatestReviews directly
	poller.processLatestReviews(ctx)

	// Verify that AddReviews was called exactly once
	if mockApp.addReviewsFuncCalled != 1 {
		t.Errorf("Expected AddReviews to be called exactly once, got %d calls", mockApp.addReviewsFuncCalled)
	}
}

// TestProcessLatestReviews_DoesNotCallAddReviewsWhenNoReviewsReturned verifies that AddReviews is not called when fetcher returns no reviews
func TestProcessLatestReviews_DoesNotCallAddReviewsWhenNoReviewsReturned(t *testing.T) {
	mockApp := &MockApp{}
	mockFetcher := &MockFetcher{
		mockedReviews: []models.AppStoreReview{}, // Empty slice - no reviews returned
	}
	poller := createTestPoller(mockApp, mockFetcher)

	ctx := context.Background()

	// Call processLatestReviews directly
	poller.processLatestReviews(ctx)

	// Verify that AddReviews was not called
	if mockApp.addReviewsFuncCalled != 0 {
		t.Errorf("Expected AddReviews to not be called when no reviews returned, got %d calls", mockApp.addReviewsFuncCalled)
	}
}

// TestProcessLatestReviews_PassesNilLatestReviewIdWhenNotAvailable verifies that nil is passed to fetchReviews when no latest review is available
func TestProcessLatestReviews_PassesNilLatestReviewIdWhenNotAvailable(t *testing.T) {
	mockApp := &MockApp{} // GetLatestReview() returns nil by default
	mockFetcher := &MockFetcher{
		mockedReviews: []models.AppStoreReview{
			{
				ID:        "test-review-1",
				Title:     "Test Review 1",
				Content:   "Test content 1",
				Author:    "Test Author 1",
				Rating:    5,
				UpdatedAt: time.Now(),
			},
		},
	}
	poller := createTestPoller(mockApp, mockFetcher)

	ctx := context.Background()

	// Call processLatestReviews directly
	poller.processLatestReviews(ctx)

	// Verify that fetchReviews was called with nil latestReviewId
	if mockFetcher.fetchReviewsCalled != 1 {
		t.Errorf("Expected fetchReviews to be called exactly once, got %d calls", mockFetcher.fetchReviewsCalled)
	}

	if mockFetcher.capturedLatestReviewId != nil {
		t.Errorf("Expected latestReviewId to be nil when no latest review is available, got %v", mockFetcher.capturedLatestReviewId)
	}

	// Verify the appID was passed correctly
	expectedAppID := "test-app-id"
	if mockFetcher.capturedAppID != expectedAppID {
		t.Errorf("Expected appID to be %s, got %s", expectedAppID, mockFetcher.capturedAppID)
	}
}

// TestProcessLatestReviews_PassesLatestReviewIdWhenAvailable verifies that the latest review ID is passed to fetchReviews when available
func TestProcessLatestReviews_PassesLatestReviewIdWhenAvailable(t *testing.T) {
	expectedLatestReviewID := "existing-review-123"
	mockApp := &MockApp{
		mockedLatestReview: &models.AppStoreReview{
			ID:        expectedLatestReviewID,
			Title:     "Existing Review",
			Content:   "Existing content",
			Author:    "Existing Author",
			Rating:    4,
			UpdatedAt: time.Now().Add(-24 * time.Hour), // 1 day ago
		},
	}
	mockFetcher := &MockFetcher{
		mockedReviews: []models.AppStoreReview{
			{
				ID:        "new-review-456",
				Title:     "New Review",
				Content:   "New content",
				Author:    "New Author",
				Rating:    5,
				UpdatedAt: time.Now(),
			},
		},
	}
	poller := createTestPoller(mockApp, mockFetcher)

	ctx := context.Background()

	// Call processLatestReviews directly
	poller.processLatestReviews(ctx)

	// Verify that fetchReviews was called with the correct latestReviewId
	if mockFetcher.fetchReviewsCalled != 1 {
		t.Errorf("Expected fetchReviews to be called exactly once, got %d calls", mockFetcher.fetchReviewsCalled)
	}

	if mockFetcher.capturedLatestReviewId == nil {
		t.Error("Expected latestReviewId to not be nil when latest review is available")
	} else if *mockFetcher.capturedLatestReviewId != expectedLatestReviewID {
		t.Errorf("Expected latestReviewId to be %s, got %s", expectedLatestReviewID, *mockFetcher.capturedLatestReviewId)
	}

	// Verify the appID was passed correctly
	expectedAppID := "test-app-id"
	if mockFetcher.capturedAppID != expectedAppID {
		t.Errorf("Expected appID to be %s, got %s", expectedAppID, mockFetcher.capturedAppID)
	}
}

// TestProcessLatestReviews_HandlesErrorGracefully verifies that the poller handles fetchReviews errors gracefully
func TestProcessLatestReviews_HandlesErrorGracefully(t *testing.T) {
	mockApp := &MockApp{}
	mockFetcher := &MockFetcher{
		mockedError: errors.New("network connection failed"),
	}
	poller := createTestPoller(mockApp, mockFetcher)

	ctx := context.Background()

	// Call processLatestReviews directly - this should not panic or crash
	poller.processLatestReviews(ctx)

	// Verify that fetchReviews was called
	if mockFetcher.fetchReviewsCalled != 1 {
		t.Errorf("Expected fetchReviews to be called exactly once, got %d calls", mockFetcher.fetchReviewsCalled)
	}

	// Verify that AddReviews was NOT called when an error occurs
	if mockApp.addReviewsFuncCalled != 0 {
		t.Errorf("Expected AddReviews to not be called when fetchReviews returns an error, got %d calls", mockApp.addReviewsFuncCalled)
	}
}
