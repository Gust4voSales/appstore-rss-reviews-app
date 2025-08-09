package models

import (
	"slices"
	"time"
)

type AppStoreReview struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	Rating    int       `json:"rating"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type AppStoreReviews []AppStoreReview

func (reviews AppStoreReviews) Sort() {
	slices.SortFunc(reviews, func(a, b AppStoreReview) int {
		if a.UpdatedAt.After(b.UpdatedAt) {
			return -1 // a comes before b (descending order)
		}
		if a.UpdatedAt.Before(b.UpdatedAt) {
			return 1 // b comes before a
		}
		return 0 // equal
	})
}
