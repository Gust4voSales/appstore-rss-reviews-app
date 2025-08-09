package appstore_reviews_poller

import (
	"context"
	"log"
	"time"

	"github.com/Gust4voSales/appstore-rss-reviews-app/server/config"
	"github.com/Gust4voSales/appstore-rss-reviews-app/server/internal/app"
	"github.com/Gust4voSales/appstore-rss-reviews-app/server/internal/models"
)

type AppStoreReviewsPoller struct {
	cfg        *config.Config
	appService *app.App
	fetcher    *Fetcher
}

func New(cfg *config.Config, appService *app.App) *AppStoreReviewsPoller {
	fetcher := NewFetcher("https://itunes.apple.com/us/rss/customerreviews")
	return &AppStoreReviewsPoller{
		cfg:        cfg,
		appService: appService,
		fetcher:    fetcher,
	}
}

func (p *AppStoreReviewsPoller) Run(ctx context.Context) {
	ticker := time.NewTicker(p.cfg.PollingInterval)
	defer ticker.Stop()

	p.processLatestReviews(ctx) // first run immediately

	for {
		select {
		case <-ctx.Done():
			log.Println("poller stopping")
			return
		case <-ticker.C:
			p.processLatestReviews(ctx)
		}
	}
}

// fetchLatestReviews fetches all reviews with pagination
func (p *AppStoreReviewsPoller) fetchLatestReviews(ctx context.Context, latestReviewId *string) ([]models.AppStoreReview, error) {
	log.Printf(" > FETCH: fetching reviews - appId: %s", p.cfg.AppID)

	reviews, err := p.fetcher.fetchReviews(ctx, p.cfg.AppID, latestReviewId)
	if err != nil {
		return nil, err
	}

	if len(reviews) > 0 { // debug helper
		log.Printf(" > FETCH: 1st review example ID: %+v", reviews[0].ID)
	}

	return reviews, nil
}

// processLatestReviews fetches and processes all latest reviews it can find that are not already in the database
func (p *AppStoreReviewsPoller) processLatestReviews(ctx context.Context) {
	log.Printf(">> PROCESS: starting processLatestReviews <<")

	var latestReviewId *string
	latestReview := p.appService.GetLatestReview()

	if latestReview != nil {
		latestReviewId = &latestReview.ID
		log.Printf(" > PROCESS: latest review saved in db: %s", latestReview.ID)
	} else {
		log.Printf(" > PROCESS: no latest review found in db, fetching all possible reviews")
	}

	reviews, err := p.fetchLatestReviews(ctx, latestReviewId)
	if err != nil {
		log.Printf(" > ❌ PROCESS: error fetching reviews: %v", err)
		// intentionally not doing error handling here, if it fails, it will be retried in the next tick
		return
	}

	if len(reviews) == 0 {
		log.Printf(" > PROCESS: no new reviews found ✅")
		return
	}

	log.Printf(" > PROCESS: found %d reviews", len(reviews))

	added := p.appService.AddReviews(reviews)
	if added > 0 {
		log.Printf(" > PROCESS: added %d new reviews ✅", added)
	} else {
		log.Printf(" > PROCESS: no new reviews added ✅")
	}
}
