package news

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TavilySearchRequest struct {
	APIKey         string   `json:"api_key"`
	Query          string   `json:"query"`
	SearchDepth    string   `json:"search_depth"`
	MaxResults     int      `json:"max_results"`
	IncludeAnswer  bool     `json:"include_answer"`
	IncludeDomains []string `json:"include_domains,omitempty"`
	TimeRange      string   `json:"time_range,omitempty"`
}

type TavilySearchResponse struct {
	Results []TavilyResult `json:"results"`
	Answer  string         `json:"answer,omitempty"`
}

type TavilyResult struct {
	Title   string  `json:"title"`
	URL     string  `json:"url"`
	Content string  `json:"content"`
	Score   float64 `json:"score"`
}

func SearchStockNews(ticker string, apiKey string) ([]string, error) {
	query := fmt.Sprintf("%s stock price surge news today reasons", ticker)

	reqBody := TavilySearchRequest{
		APIKey:        apiKey,
		Query:         query,
		SearchDepth:   "basic",
		MaxResults:    5,
		IncludeAnswer: false,
		TimeRange:     "week",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{Timeout: 15 * time.Second}

	res, err := client.Post(
		"https://api.tavily.com/search",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to call Tavily API: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tavily API returned status: %s", res.Status)
	}

	var result TavilySearchResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var headlines []string
	for _, item := range result.Results {
		newsText := fmt.Sprintf("Title: %s\nContent: %s\nURL: %s",
			item.Title,
			item.Content,
			item.URL,
		)
		headlines = append(headlines, newsText)
	}

	return headlines, nil
}
