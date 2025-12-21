// 実質的に使用されていないファイルなので、必要がないと判断すれば都合の良いタイミングで消去してもいい
package news

import (
	"net/http"
	"time"
	"encoding/json"
	"fmt"
)

type NewsResponse struct {
	Feed []NewsItem `json:"feed"`
}

type NewsItem struct {
	Title string `json:"title"`
	Summary string `json:"summary"`
	URL string `json:"url"`
	Time string `json:"time"`
}

func FetchNews(ticker string, apiKey string)([]string, error) {
	url := fmt.Sprintf(
		"https://www.alphavantage.co/query?function=NEWS_SENTIMENT&tickers=%s&sort=LATEST&limit=3&apikey=%s",
		ticker, apiKey,
	)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch news: %s", res.Status)
	}

	var result NewsResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	var headlines []string
	for _, item := range result.Feed {
		newsText := fmt.Sprintf("Title: %s\nSummary: %s", item.Title, item.Summary)
		headlines = append(headlines, newsText)
	}

	return headlines, nil
}