package repositories

import (
	"log"
	"time"

	"github.com/Gust4voSales/appstore-rss-reviews-app/server/internal/models"
)

// TODO improve later, implement file persistence

type AppReviewsRepository struct {
	Reviews []models.AppStoreReview
}

func Load() *AppReviewsRepository {
	return &AppReviewsRepository{
		Reviews: []models.AppStoreReview{},
	}
}

func (a *AppReviewsRepository) ListLatest(hours int) []models.AppStoreReview {
	cutoffTime := time.Now().UTC().Add(-time.Duration(hours) * time.Hour)

	var recentReviews []models.AppStoreReview

	// I'm considering that the reviews are already sorted by updatedAt in descending order
	// (since I'm fetching and saving them like that) so I'm choosing not sort them
	for _, review := range a.Reviews {
		if review.UpdatedAt.After(cutoffTime) {
			recentReviews = append(recentReviews, review)
		} else {
			log.Printf("review.UpdatedAt before cutoffTime: %s", review.UpdatedAt.Format(time.RFC3339))
			break
		}
	}

	return recentReviews
}

func (a *AppReviewsRepository) GetLatestReview() *models.AppStoreReview {
	if len(a.Reviews) == 0 {
		return nil
	}
	return &a.Reviews[0]
}

// hasReviewWithID checks if a review with the given ID already exists
func (a *AppReviewsRepository) hasReviewWithID(id string) bool {
	for _, review := range a.Reviews {
		if review.ID == id {
			return true
		}
	}
	return false
}

// AddBatch adds only new reviews (based on ID) to the repository
// Returns the number of new reviews added
func (a *AppReviewsRepository) AddBatch(reviews []models.AppStoreReview) int {
	initialLen := len(a.Reviews)

	// Filter out reviews that already exist
	for _, review := range reviews {
		if !a.hasReviewWithID(review.ID) {
			a.Reviews = append(a.Reviews, review)
		}
	}

	return len(a.Reviews) - initialLen
}
