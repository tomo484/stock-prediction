package japanese_stock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"stock-prediction/backend/models"
	"stock-prediction/backend/repositories"
	"strings"
	"time"
)

type TavilySearchRequest struct {
	APIKey         string   `json:"api_key"`
	Query          string   `json:"query"`
	SearchDepth    string   `json:"search_depth"`
	MaxResults     int      `json:"max_results"`
	IncludeAnswer  bool     `json:"include_answer"`
	IncludeDomains []string `json:"include_domains,omitempty"`
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

// 日本株のニュース検索を実行し、DBに保存する
func SearchJapaneseStockNews(companyName string, code string, apiKey string, repository repositories.IJapaneseStockRepository) (*models.NewsSearch, error) {
	queries := []string{
		fmt.Sprintf("%s ビジネスモデル", companyName),
		fmt.Sprintf("%s 決算短信 要約", companyName),
	}

	// NewsSearchﾚｺｰﾄﾞを作成
	newsSearch := &models.NewsSearch{
		Code:       code,
		SearchedAt: time.Now(),
	}

	// 並列で検索を実行する
	type searchResult struct {
		query  string
		result *TavilySearchResponse
		err    error
	}

	// goroutineの書き方
	resultChan := make(chan searchResult, len(queries))

	for _, query := range queries {
		go func(q string) {
			result, err := executeTavilySearch(q, apiKey)
			resultChan <- searchResult{query: q, result: result, err: err}
		}(query)
	}

	//結果を収集
	var allItems []models.NewsItem
	var combinedContents []string

	for i := 0; i < len(queries); i++ {
		sr := <-resultChan
		if sr.err != nil {
			log.Printf("Warning: Failed to execute search query %s: %v", sr.query, sr.err)
			continue
		}

		// NewsItemを作成
		for _, item := range sr.result.Results {
			newsItem := models.NewsItem{
				SearchQuery: sr.query,
				Title:       item.Title,
				URL:         item.URL,
				Content:     item.Content,
				Score:       item.Score,
			}
			allItems = append(allItems, newsItem)

			// ConbinedContent用のテキストを生成
			combinedContents = append(combinedContents,
				fmt.Sprintf("Query: %s\nTitle: %s\nContent: %s\nURL: %s\n\n",
					sr.query, item.Title, item.Content, item.URL))
		}
	}

	// 検索結果が空の場合はエラーを返す
	if len(allItems) == 0 {
		return nil, fmt.Errorf("no news items found for company %s", companyName)
	}

	// CombinedContentを作成
	newsSearch.CombinedContent = strings.Join(combinedContents, "\n---\n\n")

	// NewsSearchとNewsItemをDBに保存する（トランザクション内で）
	err := repository.CreateNewsSearchWithItems(newsSearch, allItems)
	if err != nil {
		return nil, fmt.Errorf("failed to create news search with items: %w", err)
	}

	return newsSearch, nil
}

// コードを体系的に理解したい。何を明確化すべきか？今回であれば、TavilySearchの部分の実装をどういう風に理解すればいいかわからなかった。
// それはなぜか？TavilySearchを介してデータを取得する関数は以前も実装したが、その際は一つのqueryしか使用していなかった。今回は複数queryを使用する必要があった。
// その違いによって、困惑した。でたもそれを文章にして尋ねたところCursorからしっかりしとした答えが返ってきた。
// そして、for ループの繰り返しより、並列化の方がいいのではないかというアドバイスをくれた。であれば、この方法でいいのではないだろうか。
// まずやりたいことを言語化する。それがもやがかかっている状態であれば、データ遷移を図におこしたら、必ず俺の頭であれば問題を細分化することができる。
// 問題を出した後に、そういう風なコードがないかどうかを探す。自力で探してもいいし、Cursorに探させてもいい。
// もしあれば、それを参考にしてこういうコードを書けばいいよね？っていうのをCursorに尋ねる。
// なければ、こういうことしたいんだけど、どういう風にすればいいという形でCursorに尋ねる。
// 決して、何も分からないままCursorに尋ねるのは避ける。めんどくさくても、解像度が低くてもいいから答えを持った状態でCursorにモノを尋ねる。
// 話はそこから。これでよくない？実装を進める際に毎回これを確認して、実装を進めていけば何の問題もないやん。
// では、具体例をやりましょうか。①何をしたいのか = 取得したニュース情報（models.NewsSearch, models.NewsItemをDBに保存したい。）
// ②具体例は存在するのか→あるやろ。取得した情報をDBに保存するという処理が書かれている部分を探してくるだけです。→daily_stock, companyInfoとかに存在するんじゃない？
// ①外部テーブルを参照するためのIDを途中で付与したい。その場合はどうすればいいのだろうか、NewsSerachのカラムを取得するGetメソッドが必要で、latestを取得して、それに +1をするみたいな実装の仕方をするのかな？

// executeTavilySearch Tavily APIを呼び出す内部関数
func executeTavilySearch(query string, apiKey string) (*TavilySearchResponse, error) {
	reqBody := TavilySearchRequest{
		APIKey:        apiKey,
		Query:         query,
		SearchDepth:   "basic",
		MaxResults:    5,
		IncludeAnswer: false,
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

	return &result, nil
}

func SyncJapaneseStockNews(companyName string, code string, apiKey string, repository repositories.IJapaneseStockRepository) (*models.NewsSearch, error) {
	// Tavily APIからニュースを検索し、DBに保存
	_, err := SearchJapaneseStockNews(companyName, code, apiKey, repository)
	if err != nil {
		return nil, fmt.Errorf("failed to search and save Japanese stock news: %w", err)
	}

	// 取得し、保存したニュース検索データを返す（型の整合性とデータの正確性を保証）
	newsSearch, err := repository.FindNewsByCode(code)
	if err != nil {
		return nil, fmt.Errorf("failed to find news by code: %w", err)
	}

	return newsSearch, nil
}

// SyncJapaneseStockNewsのメソッドを並列化したもの
func SyncJapaneseStockNewsParallel(companies []models.Company, apiKey string, repository repositories.IJapaneseStockRepository)([]*models.NewsSearch, error) {
	if len(companies) == 0 {
		return nil, fmt.Errorf("no companies provided")
	}

	// 結果を格納する構造体（何の為に存在するのかが正直よくわからない）
	type syncResult struct {
		company models.Company 
		newsSearch *models.NewsSearch
		err error
	}

	// 並列化でよくある格納用のmakeメソッド
	resultChan := make(chan syncResult, len(companies))

	// 各企業に対してgoroutineでSyncJapaneseStockNewsを実行
	for _, company := range companies {
		go func(c models.Company) {
			// 1社分のnews同期処理を実行
			newsSearch, err := SyncJapaneseStockNews(c.CompanyName, c.Code, apiKey, repository)
			resultChan <- syncResult{company: c, newsSearch: newsSearch, err: err,}
		}(company)
	}

	// 結果を収集
	var successResults []*models.NewsSearch
	var errors []error

	for i := 0; i < len(companies); i ++ {
		result := <-resultChan
		if result.err != nil {
			log.Printf("Warning: Failed to sync news for company %s (code: %s): %v",result.company.CompanyName, result.company.Code, result.err)
			errors = append(errors, fmt.Errorf("company %s: %w", result.company.Code, result.err))
			continue
		}
		successResults = append(successResults, result.newsSearch)
	}

	// 全ての企業で失敗した場合はエラーを返す
	if len(successResults) == 0 {
		return nil, fmt.Errorf("failed to sync news for all companies: %v", errors)
	}

	//一部成功した場合は成功した結果を返す（エラーはログに記録済み）
	if len(errors) > 0 {
		log.Printf("Info: Successfully synced news for %d/%d companies", len(successResults), len(companies))
	}

	return successResults, nil
}