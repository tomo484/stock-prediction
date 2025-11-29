package router

import (
	"stock-prediction/backend/controllers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(sc controllers.IStockController) *echo.Echo {
	e := echo.New()

	// CORS設定
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:3004"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}))

	// ログミドルウェア
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// API routes
	api := e.Group("/api")

	// Stocks routes
	stocks := api.Group("/stocks")
	stocks.GET("/latest", sc.FindLatestRanking)
	stocks.GET("/date", sc.FindDailyRanking)
	stocks.GET("/:ticker", sc.FindStock)

	// Admin routes
	admin := api.Group("/admin")
	admin.POST("/sync", sc.SyncData)

	return e
}
