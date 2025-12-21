package america_stock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"stock-prediction/backend/models"
	"stock-prediction/backend/repositories"
	"stock-prediction/backend/utils"
	"time"
)

type FMPResponse struct {
	Symbol            string  `json:"symbol"`
	Price             float64 `json:"price"`
	MarketCap         float64 `json:"marketCap"`
	Beta              float64 `json:"beta"`
	LastDividend      float64 `json:"lastDividend"`
	Range             string  `json:"range"`
	Change            float64 `json:"change"`
	ChangePercentage  float64 `json:"changePercentage"`
	Volume            int64   `json:"volume"`
	AverageVolume     int64   `json:"averageVolume"`
	CompanyName       string  `json:"companyName"`
	Currency          string  `json:"currency"`
	Cik               string  `json:"cik"`
	Isin              string  `json:"isin"`
	Cusip             string  `json:"cusip"`
	ExchangeFullName  string  `json:"exchangeFullName"`
	Exchange          string  `json:"exchange"`
	Industry          string  `json:"industry"`
	Website           string  `json:"website"`
	Description       string  `json:"description"`
	CEO               string  `json:"ceo"`
	Sector            string  `json:"sector"`
	Country           string  `json:"country"`
	FullTimeEmployees string  `json:"fullTimeEmployees"` // APIは文字列で返す
	Phone             string  `json:"phone"`
	Address           string  `json:"address"`
	City              string  `json:"city"`
	State             string  `json:"state"`
	Zip               string  `json:"zip"`
	Image             string  `json:"image"`
	IpoDate           string  `json:"ipoDate"`
	DefaultImage      bool    `json:"defaultImage"`
	IsEtf             bool    `json:"isEtf"`
	IsActivelyTrading bool    `json:"isActivelyTrading"`
	IsAdr             bool    `json:"isAdr"`
	IsFund            bool    `json:"isFund"`
}

func FetchFMPData(ticker string, apiKey string) (*FMPResponse, error) {
	url := fmt.Sprintf("https://financialmodelingprep.com/stable/profile?symbol=%s&apikey=%s", ticker, apiKey)
	client := &http.Client{Timeout: 10 * time.Second}

	res, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("データの取得に失敗しました: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("データの取得に失敗しました(Status関連で): %s", res.Status)
	}

	var result []FMPResponse
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("データのでコードに失敗しました: %w", err)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("ticker: %s のデータが見つかりませんでした", ticker)
	}

	return &result[0], nil
}

func SaveFMPDatatoDB(fmpData *FMPResponse, repo repositories.IStockRepository) error {
	// 既存のStockテーブルのデータを取得
	stock, err := repo.FindStockByTicker(fmpData.Symbol)
	if err != nil {
		return fmt.Errorf("stock %s not found in DB: %w", fmpData.Symbol, err)
	}
	// 静的データをStockテーブルに保存
	if stock.Name == "" {
		stock.Name = fmpData.CompanyName
		stock.Sector = fmpData.Sector
		stock.Industry = fmpData.Industry
		stock.Description = fmpData.Description
		stock.Website = fmpData.Website
		stock.Country = fmpData.Country
		stock.FullTimeEmployees = int(utils.ParseInt(fmpData.FullTimeEmployees))
		stock.Image = fmpData.Image
		stock.IpoDate = fmpData.IpoDate
		stock.CEO = fmpData.CEO
	}
	if err := repo.UpdateStock(stock); err != nil {
		return fmt.Errorf("failed to create/update stock %s: %w", fmpData.Symbol, err)
	}

	// 動的データをStockMetricテーブルに保存
	today := time.Now().Format("2006-01-02")

	metric := &models.StockMetric{
		StockID:       stock.ID,
		Date:          today,
		MarketCap:     fmpData.MarketCap,
		Volume:        fmpData.Volume,
		AverageVolume: fmpData.AverageVolume,
		Beta:          fmpData.Beta,
		LastDividend:  fmpData.LastDividend,
	}

	if err := repo.UpdateStockMetric(metric); err != nil {
		return fmt.Errorf("failed to update stock metric %s: %w", fmpData.Symbol, err)
	}

	return nil
}

func SyncCompanyInfo(ticker string, repo repositories.IStockRepository, apiKey string) error {
	// FMP APIから企業データを取得
	fmpData, err := FetchFMPData(ticker, apiKey)
	if err != nil {
		return fmt.Errorf("failed to fetch FMP data for %s: %w", ticker, err)
	}

	// FMP APIから取得した企業データをDBに保存
	if err := SaveFMPDatatoDB(fmpData, repo); err != nil {
		return fmt.Errorf("failed to save FMP data for %s: %w", ticker, err)
	}

	return nil
}