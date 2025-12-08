package xpost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"stock-prediction/backend/repositories"
	"time"

	"github.com/dghubble/oauth1"
)

type IXPostService interface {
	PostRanking(date string) error
	PostAnalysis(date string) error
	PostSingleAnalysis(date string, rank int) error
}

type xPostService struct {
	repository repositories.IStockRepository
	client     *http.Client
}

func NewXPostService(repo repositories.IStockRepository) IXPostService {
	// OAuth1の設定（後々必ず使用するので先にここで行っておく）
	config := oauth1.NewConfig(os.Getenv("X_API_KEY"), os.Getenv("X_POST_SECRET"))
	token := oauth1.NewToken(os.Getenv("X_ACCESS_TOKEN"), os.Getenv("X_ACCESS_TOKEN_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)

	return &xPostService{repository: repo, client: httpClient}
}

// ランキング投稿（AiAnalysis無し）
func (s *xPostService) PostRanking(date string) error {
	rankings, err := s.repository.FindDailyRanking(date)
	if err != nil {
		return fmt.Errorf("failed to find daily rankings:%w", err)
	}

	text := BuildRankingPost(date, *rankings)
	return s.postToX(text)
}

// 後者: 個別分析投稿（5件まとめて）
func (s *xPostService) PostAnalysis(date string) error {
	for rank := 1; rank <= 5; rank++ {
		if err := s.PostSingleAnalysis(date, rank); err != nil {
			// エラーログを残しつつ続行
			fmt.Printf("failed to post analysis for rank %d: %v\n", rank, err)
		}
		// レートリミット対策: 投稿間隔を空ける
		time.Sleep(5 * time.Second)
	}
	return nil
}

// 個別分析投稿（1件ずつ）
func (s *xPostService) PostSingleAnalysis(date string, rank int) error {
	ranking, err := s.repository.FindDailyRankingByDateAndRank(date, rank, "Top Gainers")
	if err != nil {
		return fmt.Errorf("failed to find daily ranking data:%w", err)
	}

	text := BuildAnalysisPost(*ranking)
	return s.postToX(text)
}

// X APIを使用しての投稿
func (s *xPostService) postToX(text string) error {
	url := "https://api.twitter.com/2/tweets"

	body := map[string]string{"text": text}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to post to x: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to post to x: status %s", res.Status)
	}

	return nil
}
