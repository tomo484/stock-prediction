package services

import (
	"net/http"
	"fmt"
	"time"
	"encoding/json"
	"stock-prediction/backend/models"
	"stock-prediction/backend/repositories"
)

type AlphaVantageResponse struct {
	Metadata           string       `json:"metadata"`
	LastUpdated        string       `json:"last_updated"`
	TopGainers         []TickerData `json:"top_gainers"`           // 今回使うのはこれ
	TopLosers          []TickerData `json:"top_losers"`            // 一応定義
	MostActivelyTraded []TickerData `json:"most_actively_traded"`  // 一応定義
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
	loc, _ := time.LoadLocation("America/New_York")
	today := time.Now().In(loc).Format("2006-01-02")

	// Top Gainersを保存
	if err := saveTickerDataToDB(alphaData.TopGainers, "Top Gainers", today, repo); err != nil {
		return fmt.Errorf("failed to save top gainers: %w", err)
	}
	
	// Top Losersを保存
	if err := saveTickerDataToDB(alphaData.TopLosers, "Top Losers", today, repo); err != nil {
		return fmt.Errorf("failed to save top losers: %w", err)
	}
	
	// Most Actively Tradedを保存
	if err := saveTickerDataToDB(alphaData.MostActivelyTraded, "Most Actively Traded", today, repo); err != nil {
		return fmt.Errorf("failed to save most actively traded: %w", err)
	}

	return nil
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