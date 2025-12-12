package main

import (
	"log"
	"os"
	"stock-prediction/backend/controllers"
	"stock-prediction/backend/db"
	"stock-prediction/backend/models"
	"stock-prediction/backend/repositories"
	"stock-prediction/backend/router"
	"stock-prediction/backend/services"
)

func main() {
	log.Println("Starting application...")

	// データベース接続
	dbConn := db.NewDB()
	defer db.CloseDB(dbConn)

	// Auto migrate: テーブルを自動的に作成・更新
	log.Println("Running database migration...")
	if err := dbConn.AutoMigrate(&models.Stock{}, &models.DailyRanking{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Successfully migrated database")

	// 依存性注入: Repository → Service → Controller
	stockRepo := repositories.NewStockRepository(dbConn)
	stockService := services.NewStockService(stockRepo)
	stockController := controllers.NewStockController(stockService, stockRepo)

	// ルーター設定
	e := router.NewRouter(stockController)

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
