package AI

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
		// Stock情報は既にPreloadされているので、直接アクセス可能
		stock := ranking.Stock

		log.Printf("Fetching news for %s...", stock.Ticker)

		// ニュースを取得
		tavilyApiKey := os.Getenv("TAVILY_API_KEY")
		headlines, err := news.SearchStockNews(stock.Ticker, tavilyApiKey)
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
