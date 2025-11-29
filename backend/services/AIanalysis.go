package services

import (
	"log"
	"os"
	"stock-prediction/backend/repositories"
	"stock-prediction/backend/services/news"
)

func PerformDailyAnalysis(repo repositories.IStockRepository) error {
	// Repository層からTop Gainersの上位5件を取得
	rankings, err := repo.FindTopRankingsByCategory("Top Gainers", 5)
	if err != nil {
		return err
	}

	log.Printf("Starting AI analysis for %d stocks...", len(*rankings))

	for _, ranking := range *rankings {
		// Repository層からStock情報を取得
		stock, err := repo.FindStockByID(ranking.StockID)
		if err != nil {
			log.Printf("Warning: Failed to find stock for ID %d: %v", ranking.StockID, err)
			continue
		}

		log.Printf("Fetching news for %s...", stock.Ticker)

		// ニュースを取得
		avApiKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
		headlines, err := news.FetchNews(stock.Ticker, avApiKey)
		if err != nil {
			log.Printf("Warning: Failed to fetch news for %s: %v", stock.Ticker, err)
			// エラーでも止まらず、ニュースなしで分析させる（Brave導入ならここで呼ぶ）
			headlines = []string{}
		}

		// AI分析を実行
		analysis, err := AnalyzeStockRise(stock.Ticker, ranking.ChangeRate, headlines)
		if err != nil {
			log.Printf("Warning: Failed to analyze %s: %v", stock.Ticker, err)
			continue
		}

		// 分析結果を更新
		ranking.AiAnalysis = analysis
		if err := repo.UpdateDailyRanking(&ranking); err != nil {
			log.Printf("Warning: Failed to update ranking for %s: %v", stock.Ticker, err)
			continue
		}

		log.Printf("Completed analysis for %s", stock.Ticker)
	}

	return nil
}
