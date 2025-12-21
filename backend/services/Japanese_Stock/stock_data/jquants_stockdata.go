package japanese_Stock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"stock-prediction/backend/models"
	"stock-prediction/backend/repositories"
	"time"
)

// DailyQuoteResponse J-Quants APIレスポンス用の型（DBモデルとは別）
type DailyQuoteResponse struct {
	Code             string  `json:"Code"`
	Date             string  `json:"Date"`
	Open             float64 `json:"Open"`
	High             float64 `json:"High"`
	Low              float64 `json:"Low"`
	Close            float64 `json:"Close"`
	Volume           float64 `json:"Volume"`
	TurnoverValue    float64 `json:"TurnoverValue"`    // 売買代金
	AdjustmentFactor float64 `json:"AdjustmentFactor"` // 調整係数
	AdjustmentOpen   float64 `json:"AdjustmentOpen"`   // 調整後始値
	AdjustmentHigh   float64 `json:"AdjustmentHigh"`   // 調整後高値
	AdjustmentLow    float64 `json:"AdjustmentLow"`    // 調整後安値
	AdjustmentClose  float64 `json:"AdjustmentClose"`  // 調整後終値
	AdjustmentVolume float64 `json:"AdjustmentVolume"` // 調整後出来高
}

type ListedDailyQuoteResponse struct {
	DailyQuotes []DailyQuoteResponse `json:"daily_quotes"`
}

func FetchJQuantsStockData(idToken string, code string, from string, to string) (*ListedDailyQuoteResponse, error) {
	url := fmt.Sprintf("https://api.jquants.com/v1/prices/daily_quotes?code=%s&from=%s&to=%s", code, from, to)
	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// IDトークンをヘッダーに付与
	req.Header.Set("Authorization", "Bearer "+idToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch stock data: status %d", resp.StatusCode)
	}

	var result ListedDailyQuoteResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// SaveJQuantsStockDataToDB 日足株価データをDBに保存する
func SaveJQuantsStockDataToDB(stockData *ListedDailyQuoteResponse, repository repositories.IJapaneseStockRepository) error {
	for _, quote := range stockData.DailyQuotes {
		// DailyQuoteResponse -> models.DailyQuoteに変換する
		dailyQuote := &models.DailyQuote{
			Code:             quote.Code,
			Date:             quote.Date,
			Open:             quote.Open,
			High:             quote.High,
			Low:              quote.Low,
			Close:            quote.Close,
			Volume:           quote.Volume,
			TurnoverValue:    quote.TurnoverValue,
			AdjustmentFactor: quote.AdjustmentFactor,
			AdjustmentOpen:   quote.AdjustmentOpen,
			AdjustmentHigh:   quote.AdjustmentHigh,
			AdjustmentLow:    quote.AdjustmentLow,
			AdjustmentClose:  quote.AdjustmentClose,
			AdjustmentVolume: quote.AdjustmentVolume,
		}

		// 既存なら更新、新規なら保存する
		err := repository.CreateOrUpdateDailyQuote(dailyQuote)
		if err != nil {
			return fmt.Errorf("failed to create/update daily quote %s %s: %w", quote.Code, quote.Date, err)
		}
	}

	return nil
}

func SyncJQuantsStockData(idToken string, code string, from string, to string, repository repositories.IJapaneseStockRepository)([]models.DailyQuote, error) {
	// JQuants APIから株価データを取得
	stockData, err := FetchJQuantsStockData(idToken, code, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JQuants stock data: %w", err)
	}

	// 株価データをDBに保存
	err = SaveJQuantsStockDataToDB(stockData, repository)
	if err != nil {
		return nil, fmt.Errorf("failed to save JQuants stock data to DB: %w", err)
	}

	// 取得し、保存した株価データを返す（型の整合性とデータの正確性を保証）
	dailyQuotes, err := repository.FindDailyQuotesByCode(code, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to find daily quotes by code: %w", err)
	}
	return dailyQuotes, nil
}