package repositories

// TODO improve later, implement file persistence

type AppReviewsRepository struct {
	Reviews []string // mocked string in memory data
}

func Load() *AppReviewsRepository {
	return &AppReviewsRepository{
		Reviews: []string{"review1 - mock", "review2 - mock", "review3 - mock"}, // loading mock data
	}
}

func (a *AppReviewsRepository) List() []string {
	return a.Reviews
}

func (a *AppReviewsRepository) AddBatch(reviews []string) {
	a.Reviews = append(a.Reviews, reviews...)
}
