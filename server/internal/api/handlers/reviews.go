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
		hours := 48 // default value

		ratingQuery := c.Query("rating")
		if ratingQuery != "" {
			parsedrating, err := strconv.Atoi(ratingQuery)
			if err != nil || !slices.Contains(validRatings, parsedrating) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rating parameter"})
				return
			}
			rating = &parsedrating
		}

		hoursQuery := c.Query("hours")
		if hoursQuery != "" {
			parsedHours, err := strconv.Atoi(hoursQuery)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hours parameter"})
				return
			}
			if parsedHours < 1 || parsedHours > 96 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Hours parameter must be between 1 and 96"})
				return
			}
			hours = parsedHours
		}

		reviews := appService.ListLatestReviews(hours, rating)
		c.JSON(http.StatusOK, struct {
			AppID     string                  `json:"appId"`
			Count     int                     `json:"count"`
			Reviews   []models.AppStoreReview `json:"reviews"`
			LastHours int                     `json:"lastHours"`
		}{
			AppID:     appService.GetAppID(),
			Count:     len(reviews),
			Reviews:   reviews,
			LastHours: hours,
		})
	}
}
