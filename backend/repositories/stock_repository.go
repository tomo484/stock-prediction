package repositories

import (
	"errors"
	"stock-prediction/backend/models"
	"gorm.io/gorm"
)

type IStockRepository interface {
	FindLatestRanking() (*[]models.DailyRanking, error)
	FindDailyRanking(date string) (*[]models.DailyRanking, error)
	FindStock(ticker string) (*[]models.DailyRanking, error)
	CreateOrUpdateStock(stock *models.Stock) error
	CreateOrUpdateDailyRanking(ranking *models.DailyRanking) error
	FindTopRankingsByCategory(category string, limit int) (*[]models.DailyRanking, error)
	FindStockByID(id uint) (*models.Stock, error)
	UpdateDailyRanking(ranking *models.DailyRanking) error
}

type stockrepository struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) IStockRepository {
	return &stockrepository{db: db}
}

func (r *stockrepository) FindLatestRanking() (*[]models.DailyRanking, error) {
	var dailyRanking []models.DailyRanking
	
	// 1. 最新の日付を取得
	var latestDate string
	dateResult := r.db.Model(&models.DailyRanking{}).
		Select("MAX(date)").
		Scan(&latestDate)
	if dateResult.Error != nil {
		return nil, dateResult.Error
	}
	if latestDate == "" {
		return nil, errors.New("no data found")
	}
	
	// 2. 最新日付のTop Gainersの1~5位を取得
	result := r.db.Preload("Stock").
		Where("date = ? AND category = ? AND rank <= ?", latestDate, "Top Gainers", 5).
		Order("rank ASC").
		Find(&dailyRanking)
	
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("no data found")
		}
		return nil, result.Error
	}
	
	return &dailyRanking, nil
}

func (r *stockrepository) FindDailyRanking(date string) (*[]models.DailyRanking, error) {
	var dailyRanking []models.DailyRanking
	// 指定日付のTop Gainersの1~5位を取得
	result := r.db.Preload("Stock").
		Where("date = ? AND category = ? AND rank <= ?", date, "Top Gainers", 5).
		Order("rank ASC").
		Find(&dailyRanking)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("no data found")
		}
		return nil, result.Error
	}
	return &dailyRanking, nil
}

func (r *stockrepository) FindStock(ticker string) (*[]models.DailyRanking, error) {
	var dailyRanking []models.DailyRanking
	// StockテーブルとJOINしてtickerで検索
	result := r.db.Preload("Stock").
		Joins("JOIN stocks ON daily_rankings.stock_id = stocks.id").
		Where("stocks.ticker = ?", ticker).
		Order("daily_rankings.created_at DESC").
		Find(&dailyRanking)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("no data found")
		}
		return nil, result.Error
	}
	return &dailyRanking, nil
}

func (r *stockrepository) CreateOrUpdateStock(stock *models.Stock) error {
	var existingStock models.Stock
	result := r.db.Where("ticker = ?", stock.Ticker).First(&existingStock)

	if result.Error == gorm.ErrRecordNotFound {
		return r.db.Create(stock).Error
	} else if result.Error != nil {
		return result.Error
	}

	stock.ID = existingStock.ID
	return nil
}

func (r *stockrepository)  CreateOrUpdateDailyRanking(ranking *models.DailyRanking) error {
	var existingRanking models.DailyRanking
	result := r.db.Where("date = ? AND stock_id = ? AND category = ?", 
	ranking.Date, ranking.StockID, ranking.Category).First(&existingRanking)

	if result.Error == gorm.ErrRecordNotFound {
		// 新規作成
		return r.db.Create(ranking).Error
	} else if result.Error != nil {
		return result.Error
	}

	ranking.ID = existingRanking.ID
	return r.db.Model(&existingRanking).Updates(ranking).Error
}

func (r *stockrepository) FindTopRankingsByCategory(category string, limit int) (*[]models.DailyRanking, error) {
	var rankings []models.DailyRanking
	result := r.db.Where("category = ? AND rank <= ?", category, limit).
		Order("rank ASC").
		Find(&rankings)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("no data found")
		}
		return nil, result.Error
	}
	
	return &rankings, nil
}

func (r *stockrepository) FindStockByID(id uint) (*models.Stock, error) {
	var stock models.Stock
	result := r.db.First(&stock, id)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("stock not found")
		}
		return nil, result.Error
	}
	
	return &stock, nil
}

func (r *stockrepository) UpdateDailyRanking(ranking *models.DailyRanking) error {
	return r.db.Save(ranking).Error
}