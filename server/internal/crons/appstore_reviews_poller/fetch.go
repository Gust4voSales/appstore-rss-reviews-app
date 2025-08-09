package appstore_reviews_poller

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/Gust4voSales/appstore-rss-reviews-app/server/internal/models"
)

// Fetcher handles fetching reviews from the App Store RSS feed
type Fetcher struct {
	client  *http.Client
	baseURL string
}

// NewFetcher creates a new Fetcher instance
func NewFetcher(baseURL string) *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
	}
}

// fetchReviews fetches all reviews with pagination support
func (f *Fetcher) fetchReviews(ctx context.Context, appID string, latestReviewId *string) ([]models.AppStoreReview, error) {
	var allReviews []models.AppStoreReview

	// Start with page 1 and continue
	page := 1
	foundLatestReview := false
	const maxPages = 10 // Safety limit to prevent infinite loops and it is also the maximum number of pages Apple allows

	for page <= maxPages {
		time.Sleep(200 * time.Millisecond) // sleep to avoid potential rate limiting
		reviews, err := f.fetchMostRecentReviewsPage(ctx, appID, page)
		if err != nil {
			return nil, fmt.Errorf("fetching page %d: %w", page, err)
		}

		// If no reviews returned, we've reached the end
		if len(reviews) == 0 {
			log.Printf(" > FETCH: no reviews found on page %d, stopping pagination", page)
			break
		}

		for _, review := range reviews {
			// if we have a latestReviewId, we stop fetching when we reach it
			if latestReviewId != nil && review.ID == *latestReviewId {
				foundLatestReview = true
				break
			}
			allReviews = append(allReviews, review)
		}

		if foundLatestReview {
			break
		}

		page++
	}

	return allReviews, nil
}

// fetchMostRecentReviewsPage fetches a specific page of reviews
func (f *Fetcher) fetchMostRecentReviewsPage(ctx context.Context, appID string, page int) ([]models.AppStoreReview, error) {
	log.Printf(" > FETCH: fetching page: %d", page)

	// Build URL with page parameter
	baseURL := fmt.Sprintf("%s/id=%s/sortBy=mostRecent/page=%d/json", f.baseURL, appID, page)
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("parsing URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	reviews, err := parseAppStoreReviews(body)
	if err != nil {
		return nil, fmt.Errorf("parsing reviews: %w", err)
	}

	return reviews, nil
}
