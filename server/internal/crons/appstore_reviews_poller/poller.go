package appstore_reviews_poller

import (
	"context"
	"log"
	"time"

	"github.com/Gust4voSales/appstore-rss-reviews-app/server/config"
	"github.com/Gust4voSales/appstore-rss-reviews-app/server/internal/app"
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

	p.doOnce(ctx) // first run immediately

	for {
		select {
		case <-ctx.Done():
			log.Println("poller stopping")
			return
		case <-ticker.C:
			p.doOnce(ctx)
		}
	}
}

func (p *AppStoreReviewsPoller) doOnce(ctx context.Context) {
	log.Println("⏲️ CRON: polling for reviews - appId: ", p.cfg.AppID)
	reviews, err := p.fetcher.FetchReviews(ctx, p.cfg.AppID)
	if err != nil {
		// in a production environment, we would have monitoring to watch for this error
		log.Printf("❌ CRON: error fetching reviews: %v", err)
		return
	}

	if len(reviews) == 0 {
		log.Println("❌ CRON: no reviews found")
		return
	}

	log.Printf("✅ CRON: found %d reviews", len(reviews))
	log.Printf("✅ CRON: 1st review example: %+v", reviews[0]) // TODO remove debug

	added := p.appService.AddReviews(reviews)
	if added > 0 {
		log.Printf("✅ CRON: added %d new reviews", added)
	} else {
		log.Println("✅ CRON: no new reviews")
	}
}
