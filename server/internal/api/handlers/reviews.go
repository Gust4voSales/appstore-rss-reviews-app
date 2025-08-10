package handlers

import (
	"net/http"
	"slices"
	"strconv"

	"github.com/Gust4voSales/appstore-rss-reviews-app/server/internal/app"
	"github.com/Gust4voSales/appstore-rss-reviews-app/server/internal/models"
	"github.com/gin-gonic/gin"
)

var validRatings = []int{1, 2, 3, 4, 5}

func ListReviews(appService *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rating *int

		ratingQuery := c.Query("rating")
		if ratingQuery != "" {
			parsedrating, err := strconv.Atoi(ratingQuery)
			if err != nil || !slices.Contains(validRatings, parsedrating) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rating parameter"})
				return
			}
			rating = &parsedrating
		}

		reviews := appService.ListLatestReviews(96, rating)
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
