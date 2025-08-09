package appstore_reviews_poller

import (
	"context"
	"fmt"
	"io"
	"net/http"
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

// FetchReviews fetches the latest reviews for a given app ID
func (f *Fetcher) FetchReviews(ctx context.Context, appID string) ([]models.AppStoreReview, error) {
	url := fmt.Sprintf("%s/id=%s/sortBy=mostRecent/json", f.baseURL, appID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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
