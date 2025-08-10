package appstore_reviews_poller

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestFetchMostRecentReviewsPage_Success verifies that fetching a single page of reviews works correctly with valid JSON response
func TestFetchMostRecentReviewsPage_Success(t *testing.T) {
	// Mock response data
	mockResponse := `{
		"feed": {
			"entry": [
				{
					"id": {"label": "review-1"},
					"title": {"label": "Great app!"},
					"content": {"label": "I love this app"},
					"author": {"name": {"label": "John Doe"}},
					"im:rating": {"label": "5"},
					"updated": {"label": "2023-12-07T10:30:00-07:00"}
				},
				{
					"id": {"label": "review-2"},
					"title": {"label": "Good app"},
					"content": {"label": "Pretty decent"},
					"author": {"name": {"label": "Jane Smith"}},
					"im:rating": {"label": "4"},
					"updated": {"label": "2023-12-06T15:20:00-07:00"}
				}
			]
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the URL format
		expectedPath := "/id=123456789/sortBy=mostRecent/page=1/json"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	fetcher := NewFetcher(server.URL)
	ctx := context.Background()

	reviews, err := fetcher.fetchMostRecentReviewsPage(ctx, "123456789", 1)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(reviews) != 2 {
		t.Errorf("Expected 2 reviews, got %d", len(reviews))
	}

	// Verify first review
	if reviews[0].ID != "review-1" {
		t.Errorf("Expected first review ID 'review-1', got '%s'", reviews[0].ID)
	}
	if reviews[0].Title != "Great app!" {
		t.Errorf("Expected first review title 'Great app!', got '%s'", reviews[0].Title)
	}
	if reviews[0].Rating != 5 {
		t.Errorf("Expected first review rating 5, got %d", reviews[0].Rating)
	}
}

// TestFetchMostRecentReviewsPage_HTTPError verifies that HTTP errors are properly handled and returned
func TestFetchMostRecentReviewsPage_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	fetcher := NewFetcher(server.URL)
	ctx := context.Background()

	reviews, err := fetcher.fetchMostRecentReviewsPage(ctx, "123456789", 1)
	if err == nil {
		t.Fatal("Expected error for HTTP 500, got nil")
	}

	if reviews != nil {
		t.Error("Expected nil reviews on error")
	}

	if !strings.Contains(err.Error(), "unexpected status code 500") {
		t.Errorf("Expected error message to contain status code, got: %v", err)
	}
}

// TestFetchMostRecentReviewsPage_InvalidJSON verifies that invalid JSON responses are properly handled
func TestFetchMostRecentReviewsPage_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	fetcher := NewFetcher(server.URL)
	ctx := context.Background()

	reviews, err := fetcher.fetchMostRecentReviewsPage(ctx, "123456789", 1)
	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}

	if reviews != nil {
		t.Error("Expected nil reviews on error")
	}

	if !strings.Contains(err.Error(), "parsing reviews") {
		t.Errorf("Expected error message to contain 'parsing reviews', got: %v", err)
	}
}

// TestFetchMostRecentReviewsPage_Timeout verifies that timeouts are properly handled
func TestFetchMostRecentReviewsPage_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"feed": {"entry": []}}`))
	}))
	defer server.Close()

	fetcher := NewFetcher(server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	reviews, err := fetcher.fetchMostRecentReviewsPage(ctx, "123456789", 1)
	if err == nil {
		t.Fatal("Expected error for context cancellation, got nil")
	}

	if reviews != nil {
		t.Error("Expected nil reviews on error")
	}
}

// TestFetchReviews_SinglePage verifies that fetching reviews stops when only one page of data is available
func TestFetchReviews_SinglePage(t *testing.T) {
	mockResponse := `{
		"feed": {
			"entry": [
				{
					"id": {"label": "review-1"},
					"title": {"label": "Great app!"},
					"content": {"label": "I love this app"},
					"author": {"name": {"label": "John Doe"}},
					"im:rating": {"label": "5"},
					"updated": {"label": "2023-12-07T10:30:00-07:00"}
				}
			]
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "page=1") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(mockResponse))
		} else {
			// Return empty for subsequent pages
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"feed": {"entry": []}}`))
		}
	}))
	defer server.Close()

	fetcher := NewFetcher(server.URL)
	ctx := context.Background()

	reviews, err := fetcher.fetchReviews(ctx, "123456789", nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(reviews) != 1 {
		t.Errorf("Expected 1 review, got %d", len(reviews))
	}

	if reviews[0].ID != "review-1" {
		t.Errorf("Expected review ID 'review-1', got '%s'", reviews[0].ID)
	}
}

// TestFetchReviews_MultiplePages verifies that fetching reviews correctly aggregates data from multiple pages
func TestFetchReviews_MultiplePages(t *testing.T) {
	responses := map[string]string{
		"page=1": `{
			"feed": {
				"entry": [
					{
						"id": {"label": "review-1"},
						"title": {"label": "Great app!"},
						"content": {"label": "I love this app"},
						"author": {"name": {"label": "John Doe"}},
						"im:rating": {"label": "5"},
						"updated": {"label": "2023-12-07T10:30:00-07:00"}
					}
				]
			}
		}`,
		"page=2": `{
			"feed": {
				"entry": [
					{
						"id": {"label": "review-2"},
						"title": {"label": "Good app"},
						"content": {"label": "Pretty decent"},
						"author": {"name": {"label": "Jane Smith"}},
						"im:rating": {"label": "4"},
						"updated": {"label": "2023-12-06T15:20:00-07:00"}
					}
				]
			}
		}`,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		for pageParam, response := range responses {
			if strings.Contains(path, pageParam) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(response))
				return
			}
		}
		// Return empty for other pages
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"feed": {"entry": []}}`))
	}))
	defer server.Close()

	fetcher := NewFetcher(server.URL)
	ctx := context.Background()

	reviews, err := fetcher.fetchReviews(ctx, "123456789", nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(reviews) != 2 {
		t.Errorf("Expected 2 reviews, got %d", len(reviews))
	}

	// Verify reviews are in order
	if reviews[0].ID != "review-1" {
		t.Errorf("Expected first review ID 'review-1', got '%s'", reviews[0].ID)
	}
	if reviews[1].ID != "review-2" {
		t.Errorf("Expected second review ID 'review-2', got '%s'", reviews[1].ID)
	}
}

// TestFetchReviews_WithLatestReviewId verifies that fetching stops at the specified latest review ID, returning only newer reviews
func TestFetchReviews_WithLatestReviewId(t *testing.T) {
	mockResponse := `{
		"feed": {
			"entry": [
				{
					"id": {"label": "review-new"},
					"title": {"label": "New review"},
					"content": {"label": "This is new"},
					"author": {"name": {"label": "New User"}},
					"im:rating": {"label": "5"},
					"updated": {"label": "2023-12-08T10:30:00-07:00"}
				},
				{
					"id": {"label": "review-latest"},
					"title": {"label": "Latest review"},
					"content": {"label": "This is the latest we had"},
					"author": {"name": {"label": "Latest User"}},
					"im:rating": {"label": "4"},
					"updated": {"label": "2023-12-07T10:30:00-07:00"}
				},
				{
					"id": {"label": "review-old"},
					"title": {"label": "Old review"},
					"content": {"label": "This is old"},
					"author": {"name": {"label": "Old User"}},
					"im:rating": {"label": "3"},
					"updated": {"label": "2023-12-06T10:30:00-07:00"}
				}
			]
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	fetcher := NewFetcher(server.URL)
	ctx := context.Background()

	latestReviewId := "review-latest"
	reviews, err := fetcher.fetchReviews(ctx, "123456789", &latestReviewId)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should only get reviews before the latest review ID
	if len(reviews) != 1 {
		t.Errorf("Expected 1 review, got %d", len(reviews))
	}

	if reviews[0].ID != "review-new" {
		t.Errorf("Expected review ID 'review-new', got '%s'", reviews[0].ID)
	}
}

// TestFetchReviews_EmptyResponse verifies that empty API responses are handled gracefully
func TestFetchReviews_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"feed": {}}`))
	}))
	defer server.Close()

	fetcher := NewFetcher(server.URL)
	ctx := context.Background()

	reviews, err := fetcher.fetchReviews(ctx, "123456789", nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(reviews) != 0 {
		t.Errorf("Expected 0 reviews, got %d", len(reviews))
	}
}

// TestFetchReviews_MaxPagesLimit verifies that fetching stops after the maximum number of pages
func TestFetchReviews_MaxPagesLimit(t *testing.T) {
	// Create a server that always returns one review per page
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract page number from URL
		path := r.URL.Path
		var pageNum string
		if strings.Contains(path, "page=") {
			parts := strings.Split(path, "page=")
			if len(parts) > 1 {
				pageNum = strings.Split(parts[1], "/")[0]
			}
		}

		mockResponse := fmt.Sprintf(`{
			"feed": {
				"entry": [
					{
						"id": {"label": "review-page-%s"},
						"title": {"label": "Review from page %s"},
						"content": {"label": "Content from page %s"},
						"author": {"name": {"label": "User %s"}},
						"im:rating": {"label": "5"},
						"updated": {"label": "2023-12-07T10:30:00-07:00"}
					}
				]
			}
		}`, pageNum, pageNum, pageNum, pageNum)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	fetcher := NewFetcher(server.URL)
	ctx := context.Background()

	reviews, err := fetcher.fetchReviews(ctx, "123456789", nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should stop at maxPages (10)
	if len(reviews) != 10 {
		t.Errorf("Expected 10 reviews (maxPages limit), got %d", len(reviews))
	}
}
