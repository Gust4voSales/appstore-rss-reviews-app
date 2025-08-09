package app

import "github.com/Gust4voSales/appstore-rss-reviews-app/server/internal/repositories"

type App struct {
	repo *repositories.AppReviewsRepository
}

func New(repo *repositories.AppReviewsRepository) *App {
	return &App{repo: repo}
}

func (a *App) ListReviews() []string {
	return a.repo.List()
}

func (a *App) AddReviews(reviews []string) {
	a.repo.AddBatch(reviews)
}
