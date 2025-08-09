package repositories

import "github.com/Gust4voSales/appstore-rss-reviews-app/server/internal/models"

// TODO improve later, implement file persistence

type AppReviewsRepository struct {
	Reviews []models.AppStoreReview
}

func Load() *AppReviewsRepository {
	return &AppReviewsRepository{
		Reviews: []models.AppStoreReview{},
	}
}

func (a *AppReviewsRepository) List() []models.AppStoreReview {
	return a.Reviews
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
