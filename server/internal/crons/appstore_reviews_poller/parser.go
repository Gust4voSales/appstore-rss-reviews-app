package appstore_reviews_poller

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/Gust4voSales/appstore-rss-reviews-app/server/internal/models"
)

// AppStoreRSSResponse represents the structure of the App Store RSS feed json response
type AppStoreRSSResponse struct {
	Feed struct {
		Entry []struct {
			ID struct {
				Label string `json:"label"`
			} `json:"id"`
			Title struct {
				Label string `json:"label"`
			} `json:"title"`
			Content struct {
				Label string `json:"label"`
			} `json:"content"`
			Author struct {
				Name struct {
					Label string `json:"label"`
				} `json:"name"`
			} `json:"author"`
			Rating struct {
				Label string `json:"label"`
			} `json:"im:rating"`
			Updated struct {
				Label string `json:"label"`
			} `json:"updated"`
		} `json:"entry"`
	} `json:"feed"`
}

// ParseAppStoreReviews parses the raw JSON response from App Store RSS feed into AppStoreReview models
func parseAppStoreReviews(data []byte) ([]models.AppStoreReview, error) {
	var response AppStoreRSSResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("unmarshaling app store response: %w", err)
	}

	reviews := make([]models.AppStoreReview, 0, len(response.Feed.Entry))
	for _, entry := range response.Feed.Entry {
		rating, err := strconv.Atoi(entry.Rating.Label)
		if err != nil {
			rating = 0
			log.Printf("warning: parsing rating (id: %s) err: %v", entry.ID.Label, err)
		}

		review := models.AppStoreReview{
			ID:        entry.ID.Label,
			Title:     entry.Title.Label,
			Content:   entry.Content.Label,
			Author:    entry.Author.Name.Label,
			Rating:    rating,
			UpdatedAt: entry.Updated.Label,
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}
