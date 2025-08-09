package app

import (
	"github.com/Gust4voSales/appstore-rss-reviews-app/server/config"
	"github.com/Gust4voSales/appstore-rss-reviews-app/server/internal/models"
	"github.com/Gust4voSales/appstore-rss-reviews-app/server/internal/repositories"
)

type App struct {
	repo *repositories.AppReviewsRepository
	cfg  *config.Config
}

func New(repo *repositories.AppReviewsRepository, cfg *config.Config) *App {
	return &App{repo: repo, cfg: cfg}
}

func (a *App) ListReviews() []models.AppStoreReview {
	return a.repo.List()
}

// AddReviews adds new reviews to the repository and returns the number of reviews added
func (a *App) AddReviews(reviews []models.AppStoreReview) int {
	return a.repo.AddBatch(reviews)
}

func (a *App) GetAppID() string {
	return a.cfg.AppID
}
