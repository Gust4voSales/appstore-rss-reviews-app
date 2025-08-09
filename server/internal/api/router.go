package api

import (
	"github.com/Gust4voSales/appstore-rss-reviews-app/server/internal/api/handlers"
	"github.com/Gust4voSales/appstore-rss-reviews-app/server/internal/app"
	"github.com/gin-gonic/gin"
)

func NewRouter(appService *app.App) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", handlers.Health)
	r.GET("/reviews/96h", handlers.ListReviews(appService))

	return r
}
