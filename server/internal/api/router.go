package api

import (
	"github.com/Gust4voSales/appstore-rss-reviews-app/server/internal/api/handlers"
	"github.com/Gust4voSales/appstore-rss-reviews-app/server/internal/app"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(appService *app.App) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	// cors config
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	r.Use(cors.New(corsConfig))

	r.GET("/health", handlers.Health)
	r.GET("/reviews", handlers.ListReviews(appService))

	return r
}
