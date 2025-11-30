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
	apiKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("ALPHA_VANTAGE_API_KEY is not set")
	}

	alphadata, err := FetchAlphaVantageData(apiKey)
	if err != nil {
		return fmt.Errorf("failed to fetch Alpha Vantage data: %w", err)
	}

	if err := SaveAlphaVantageDatatoDB(alphadata, s.repository); err != nil {
		return fmt.Errorf("failed to save data to DB: %w", err)
	}

	if err := AI.PerformDailyAnalysis(s.repository); err != nil {
		return fmt.Errorf("failed to perform daily analysis: %w", err)
	}

	return nil
}