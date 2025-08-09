package handlers

import (
	"net/http"

	"github.com/Gust4voSales/appstore-rss-reviews-app/server/internal/app"
	"github.com/Gust4voSales/appstore-rss-reviews-app/server/internal/models"
	"github.com/gin-gonic/gin"
)

func ListReviews(appService *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		reviews := appService.ListLatestReviews(96)
		c.JSON(http.StatusOK, struct {
			AppID   string                  `json:"appId"`
			Count   int                     `json:"count"`
			Reviews []models.AppStoreReview `json:"reviews"`
		}{
			AppID:   appService.GetAppID(),
			Count:   len(reviews),
			Reviews: reviews,
		})
	}
}
