package services

import (
	"fmt"
	"os"
	"stock-prediction/backend/models"
	"stock-prediction/backend/repositories"
	AI "stock-prediction/backend/services/AI"
)

type IStockService interface {
	FindLatestRanking() (*[]models.DailyRanking, error)
	FindDailyRanking(date string) (*[]models.DailyRanking, error)
	FindStock(ticker string) (*[]models.DailyRanking, error)
	SyncData() error
}

type stockservice struct {
	repository repositories.IStockRepository
}

func NewStockService(repository repositories.IStockRepository) IStockService {
	return &stockservice{repository: repository}
}

func (s *stockservice) FindLatestRanking() (*[]models.DailyRanking, error) {
	return s.repository.FindLatestRanking()
}

func (s *stockservice) FindDailyRanking(date string) (*[]models.DailyRanking, error) {
	return s.repository.FindDailyRanking(date)
}

func (s *stockservice) FindStock(ticker string) (*[]models.DailyRanking, error) {
	return s.repository.FindStock(ticker)
}

func (s *stockservice) SyncData() error {
	AlphaVantageApiKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
	if AlphaVantageApiKey == "" {
		return fmt.Errorf("ALPHA_VANTAGE_API_KEY is not set")
	}

	FmpApiKey := os.Getenv("FMP_API_KEY")
	if FmpApiKey == "" {
		return fmt.Errorf("FMP_API_KEY is not set")
	}

    // Alpha Vantage APIからデータを取得
	alphadata, err := FetchAlphaVantageData(AlphaVantageApiKey)
	if err != nil {
		return fmt.Errorf("failed to fetch Alpha Vantage data: %w", err)
	}

	// Alpha Vantage APIから取得したデータをDBに保存
	if err := SaveAlphaVantageDatatoDB(alphadata, s.repository); err != nil {
		return fmt.Errorf("failed to save data to DB: %w", err)
	}

	//Top Gainersの企業情報を更新（静的情報は空の場合のみ、動的情報は常に更新させる）
	for _, tickerData := range alphadata.TopGainers {
		if err := SyncCompanyInfo(tickerData.Ticker, s.repository, FmpApiKey); err != nil {
			// エラーが発生してもログに記録するのみで全体は中断しない
			fmt.Printf("failed to sync company info for %s: %w", tickerData.Ticker, err)
		}
	}

	// AI分析を実行
	if err := AI.PerformDailyAnalysis(s.repository); err != nil {
		return fmt.Errorf("failed to perform daily analysis: %w", err)
	}

	return nil
}