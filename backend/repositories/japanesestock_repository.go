package repositories

import (
	"errors"
	"stock-prediction/backend/models"

	"gorm.io/gorm"
)

type IJapaneseStockRepository interface {
	CreateOrUpdateCompany(company *models.Company) error
	CreateOrUpdateDailyQuote(dailyQuote *models.DailyQuote) error
	CreateOrUpdateFinancialStatement(financialStatement *models.FinancialStatement) error
	CreateNewsSearchWithItems(newsSearch *models.NewsSearch, items []models.NewsItem) error
	FindNewsByCode(code string) (*models.NewsSearch, error)
	FindDailyQuotesByCode(code string, fromDate string, toDate string) ([]models.DailyQuote, error)
	FindFinancialStatementsByCode(code string) ([]models.FinancialStatement, error)
	CreateOrUpdateAnalysisResult(analysisResult *models.AnalysisResult) error
	CreateOrUpdateSectorAnalysisResult(sectorAnalysisResult *models.SectorAnalysisResult) error
}

type japanesestockrepository struct {
	db *gorm.DB
}

func NewJapaneseStockRepository(db *gorm.DB) IJapaneseStockRepository {
	return &japanesestockrepository{db: db}
}

func (r *japanesestockrepository) CreateOrUpdateCompany(company *models.Company) error {
	var existingCompany models.Company
	result := r.db.Where("code = ?", company.Code).First(&existingCompany)

	if result.Error == gorm.ErrRecordNotFound {
		return r.db.Create(company).Error
	} else if result.Error != nil {
		return result.Error
	}

	company.ID = existingCompany.ID
	return r.db.Model(&existingCompany).Updates(company).Error
}

func (r *japanesestockrepository) CreateOrUpdateDailyQuote(dailyQuote *models.DailyQuote) error {
	var existingDailyQuote models.DailyQuote
	result := r.db.Where("code = ? AND date = ?", dailyQuote.Code, dailyQuote.Date).First(&existingDailyQuote)

	if result.Error == gorm.ErrRecordNotFound {
		return r.db.Create(dailyQuote).Error
	} else if result.Error != nil {
		return result.Error
	}

	dailyQuote.ID = existingDailyQuote.ID
	return r.db.Model(&existingDailyQuote).Updates(dailyQuote).Error
}

func (r *japanesestockrepository) CreateOrUpdateFinancialStatement(financialStatement *models.FinancialStatement) error {
	var existingFinancialStatement models.FinancialStatement
	result := r.db.Where("disclosure_number = ?", financialStatement.DisclosureNumber).First(&existingFinancialStatement)

	if result.Error == gorm.ErrRecordNotFound {
		return r.db.Create(financialStatement).Error
	} else if result.Error != nil {
		return result.Error
	}

	financialStatement.ID = existingFinancialStatement.ID
	return r.db.Model(&existingFinancialStatement).Updates(financialStatement).Error
}

func (r *japanesestockrepository) CreateNewsSearchWithItems(newsSearch *models.NewsSearch, items []models.NewsItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. NewsSearchを作成
		err := tx.Create(newsSearch).Error
		if err != nil {
			return err
		}

		// 2. NewsItemにNewsSearchIDを設定して一括作成
		for i := range items {
			items[i].NewsSearchID = newsSearch.ID
		}

		if len(items) > 0 {
			err := tx.Create(&items).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *japanesestockrepository) FindNewsByCode(code string) (*models.NewsSearch, error) {
	var newsSearch models.NewsSearch
	result := r.db.Where("code = ?", code).Order("searched_at DESC").Preload("Items").First(&newsSearch)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("news search not found")
		}
		return nil, result.Error
	}
	return &newsSearch, nil
}

func (r *japanesestockrepository) FindDailyQuotesByCode(code string, fromDate string, toDate string) ([]models.DailyQuote, error) {
	var dailyQuotes []models.DailyQuote
	query := r.db.Where("code = ?", code)

	if fromDate != "" {
		query = query.Where("date >= ?", fromDate)
	}
	if toDate != "" {
		query = query.Where("date <= ?", toDate)
	}

	result := query.Order("date ASC").Find(&dailyQuotes)
	if result.Error != nil {
		return nil, result.Error
	}
	return dailyQuotes, nil
}

func (r *japanesestockrepository) FindFinancialStatementsByCode(code string) ([]models.FinancialStatement, error) {
	var financialStatements []models.FinancialStatement
	result := r.db.Where("code = ?", code).Order("current_fiscal_year_end_date DESC").Find(&financialStatements)
	if result.Error != nil {
		return nil, result.Error
	}
	return financialStatements, nil
}

func (r *japanesestockrepository) CreateOrUpdateAnalysisResult(analysisResult *models.AnalysisResult) error {
	var existingResult models.AnalysisResult
	// CodeとAnalyzedAtの組み合わせで検索（同じ分析セッションを識別）
	result := r.db.Where("code = ? AND analyzed_at = ?", analysisResult.Code, analysisResult.AnalyzedAt).First(&existingResult)

	if result.Error == gorm.ErrRecordNotFound {
		// 新規作成（Phase 1の結果を保存）
		return r.db.Create(analysisResult).Error
	} else if result.Error != nil {
		return result.Error
	}

	// 既存レコードを更新（Phase 2の結果を追加）
	analysisResult.ID = existingResult.ID
	return r.db.Model(&existingResult).Updates(analysisResult).Error
}

func (r *japanesestockrepository) CreateOrUpdateSectorAnalysisResult(sectorAnalysisResult *models.SectorAnalysisResult) error {
	var existingResult models.SectorAnalysisResult
	// SectorCodeとAnalyzedAtの組み合わせで検索
	result := r.db.Where("sector_code = ? AND analyzed_at = ?",
		sectorAnalysisResult.SectorCode, sectorAnalysisResult.AnalyzedAt).First(&existingResult)

	if result.Error == gorm.ErrRecordNotFound {
		return r.db.Create(sectorAnalysisResult).Error
	} else if result.Error != nil {
		return result.Error
	}

	sectorAnalysisResult.ID = existingResult.ID
	return r.db.Model(&existingResult).Updates(sectorAnalysisResult).Error
}
