package handlers

import (
	"net/http"

	"github.com/Gust4voSales/appstore-rss-reviews-app/server/internal/app"
	"github.com/gin-gonic/gin"
)

func ListReviews(appService *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		reviews := appService.ListReviews()
		c.JSON(http.StatusOK, reviews)
	}
}

// TODO remove later, just for testing purposes
func AddReviews(appService *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		appService.AddReviews([]string{"another review - mock", "another review 2 - mock"})
		c.JSON(http.StatusOK, gin.H{"message": "Reviews added successfully"})
	}
}
