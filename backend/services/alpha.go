package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"stock-prediction/backend/models"
	"stock-prediction/backend/repositories"
	"strings"
	"time"
)

type AlphaVantageResponse struct {
	Metadata           string       `json:"metadata"`
	LastUpdated        string       `json:"last_updated"`
	TopGainers         []TickerData `json:"top_gainers"`          // 今回使うのはこれ
	TopLosers          []TickerData `json:"top_losers"`           // 一応定義
	MostActivelyTraded []TickerData `json:"most_actively_traded"` // 一応定義
}

type TickerData struct {
	Ticker           string `json:"ticker"`            // "AEHL"
	Price            string `json:"price"`             // "2.53"
	ChangeAmount     string `json:"change_amount"`     // "1.32"
	ChangePercentage string `json:"change_percentage"` // "109.0909%" (最後に%が付いている点に注意)
}

func FetchAlphaVantageData(apiKey string) (*AlphaVantageResponse, error) {
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=TOP_GAINERS_LOSERS&apikey=%s", apiKey)

	client := &http.Client{Timeout: 10 * time.Second}

	res, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: %s", res.Status)
	}

	var result AlphaVantageResponse
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode data: %w", err)
	}

	return &result, nil
}

func SaveAlphaVantageDatatoDB(alphaData *AlphaVantageResponse, repo repositories.IStockRepository) error {
	// Alpha Vantage APIのlast_updatedから日付を抽出
	date, err := extractDateFromLastUpdated(alphaData.LastUpdated)
	if err != nil {
		return fmt.Errorf("failed to extract date from last_updated: %w", err)
	}

	// Top Gainersを保存
	if err := saveTickerDataToDB(alphaData.TopGainers, "Top Gainers", date, repo); err != nil {
		return fmt.Errorf("failed to save top gainers: %w", err)
	}

	// Top Losersを保存
	if err := saveTickerDataToDB(alphaData.TopLosers, "Top Losers", date, repo); err != nil {
		return fmt.Errorf("failed to save top losers: %w", err)
	}

	// Most Actively Tradedを保存
	if err := saveTickerDataToDB(alphaData.MostActivelyTraded, "Most Actively Traded", date, repo); err != nil {
		return fmt.Errorf("failed to save most actively traded: %w", err)
	}

	return nil
}

// extractDateFromLastUpdated はAlpha Vantage APIのlast_updatedフィールドから日付を抽出します
// フォーマット例: "2025-11-28 16:15:59 US/Eastern"
func extractDateFromLastUpdated(lastUpdated string) (string, error) {
	// スペースで分割して最初の部分（日付）を取得
	parts := strings.Fields(lastUpdated)
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid last_updated format: %s", lastUpdated)
	}

	dateStr := parts[0] // "2025-11-28"

	// 日付フォーマットの検証（YYYY-MM-DD）
	_, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", fmt.Errorf("invalid date format in last_updated '%s': %w", lastUpdated, err)
	}

	return dateStr, nil
}

func saveTickerDataToDB(tickerDataList []TickerData, category string, date string, repo repositories.IStockRepository) error {
	for rank, tickerData := range tickerDataList {
		stock := &models.Stock{
			Ticker:   tickerData.Ticker,
			Name:     "", // Alpha Vantageから取得できないので空文字
			Sector:   "", // 同上
			Industry: "", // 同上
		}

		if err := repo.CreateOrUpdateStock(stock); err != nil {
			return fmt.Errorf("failed to create/update stock %s: %w", tickerData.Ticker, err)
		}

		ranking := &models.DailyRanking{
			StockID:      stock.ID,
			Date:         date,
			Rank:         rank + 1, // 1位から始まる
			Category:     category,
			ChangeAmount: ParseFloat(tickerData.ChangeAmount),
			ChangeRate:   ParsePercentage(tickerData.ChangePercentage),
			Price:        ParseFloat(tickerData.Price),
			NewsSummary:  "", // 後で設定
			AiAnalysis:   "", // 後で設定
		}

		if err := repo.CreateOrUpdateDailyRanking(ranking); err != nil {
			return fmt.Errorf("failed to create/update daily ranking for %s: %w", tickerData.Ticker, err)
		}
	}
	return nil
}
