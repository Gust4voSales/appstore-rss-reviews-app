package repositories

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Gust4voSales/appstore-rss-reviews-app/server/internal/models"
)

type AppReviewsRepository struct {
	Reviews         models.AppStoreReviews
	StorageFilePath string
}

func Load(storageFilePath string) *AppReviewsRepository {
	repo := &AppReviewsRepository{
		Reviews:         models.AppStoreReviews{},
		StorageFilePath: storageFilePath,
	}

	// Load existing data from file
	if err := repo.loadFromFile(); err != nil {
		log.Printf("Error loading reviews from file: %v", err)
		log.Printf("Starting with empty reviews list")
	} else {
		log.Printf("Loaded %d reviews from storage file: %s", len(repo.Reviews), storageFilePath)
	}

	return repo
}

func (a *AppReviewsRepository) ListLatest(hours int, query ReviewFilter) models.AppStoreReviews {
	cutoffTime := time.Now().UTC().Add(-time.Duration(hours) * time.Hour)

	var recentReviews models.AppStoreReviews

	// I'm considering that the reviews are already sorted by updatedAt in descending order
	// (since I'm fetching and saving them like that) so I'm choosing not sort them
	for _, review := range a.Reviews {
		if review.UpdatedAt.After(cutoffTime) {
			if query.Rating == nil || review.Rating == *query.Rating {
				recentReviews = append(recentReviews, review)
			}
		} else {
			// log.Printf("review.UpdatedAt before cutoffTime: %s", review.UpdatedAt.Format(time.RFC3339))
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
func (a *AppReviewsRepository) AddBatch(reviews models.AppStoreReviews) (int, error) {
	initialLen := len(a.Reviews)

	// Filter out reviews that already exist
	for _, review := range reviews {
		if !a.hasReviewWithID(review.ID) {
			a.Reviews = append(a.Reviews, review)
		}
	}

	// Sort by updatedAt in descending order
	a.Reviews.Sort()

	newReviewsCount := len(a.Reviews) - initialLen

	// Persist to file if new reviews were added
	if newReviewsCount > 0 {
		if err := a.saveToFile(); err != nil {
			return 0, fmt.Errorf("error saving reviews to file: %v", err)
		} else {
			log.Printf("Saved %d new reviews to storage file", newReviewsCount)
		}
	}

	return newReviewsCount, nil
}

// saveToFile persists the reviews to the JSON file
func (a *AppReviewsRepository) saveToFile() error {
	if a.StorageFilePath == "" {
		return nil // No file path configured, skip persistence
	}

	dir := filepath.Dir(a.StorageFilePath)
	if err := os.MkdirAll(dir, os.ModeDir); err != nil {
		return err
	}

	data, err := json.MarshalIndent(a.Reviews, "", "	")
	if err != nil {
		return err
	}

	return os.WriteFile(a.StorageFilePath, data, os.ModePerm)
}

// loadFromFile loads reviews from the JSON file
func (a *AppReviewsRepository) loadFromFile() error {
	if a.StorageFilePath == "" {
		return nil // No file path configured, skip loading
	}

	// Check if file exists
	if _, err := os.Stat(a.StorageFilePath); os.IsNotExist(err) {
		log.Printf("Storage file does not exist: %s", a.StorageFilePath)
		return nil // File doesn't exist, start with empty slice
	}

	data, err := os.ReadFile(a.StorageFilePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &a.Reviews)
}
