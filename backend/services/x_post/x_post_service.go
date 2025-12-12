package xpost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	// ===== 一時的なデバッグコード（削除予定） =====
	fmt.Printf("🔍 [DEBUG] X API環境変数の確認開始\n")

	xApiKey := os.Getenv("X_API_KEY")
	xPostSecret := os.Getenv("X_POST_SECRET")
	xAccessToken := os.Getenv("X_ACCESS_TOKEN")
	xAccessTokenSecret := os.Getenv("X_ACCESS_TOKEN_SECRET")

	fmt.Printf("  X_API_KEY: %s (長さ: %d)\n",
		maskValue(xApiKey), len(xApiKey))
	fmt.Printf("  X_POST_SECRET: %s (長さ: %d)\n",
		maskValue(xPostSecret), len(xPostSecret))
	fmt.Printf("  X_ACCESS_TOKEN: %s (長さ: %d)\n",
		maskValue(xAccessToken), len(xAccessToken))
	fmt.Printf("  X_ACCESS_TOKEN_SECRET: %s (長さ: %d)\n",
		maskValue(xAccessTokenSecret), len(xAccessTokenSecret))

	// 未設定の環境変数をチェック
	missing := []string{}
	if xApiKey == "" {
		missing = append(missing, "X_API_KEY")
	}
	if xPostSecret == "" {
		missing = append(missing, "X_POST_SECRET")
	}
	if xAccessToken == "" {
		missing = append(missing, "X_ACCESS_TOKEN")
	}
	if xAccessTokenSecret == "" {
		missing = append(missing, "X_ACCESS_TOKEN_SECRET")
	}

	if len(missing) > 0 {
		fmt.Printf("⚠️  [DEBUG] 未設定の環境変数: %v\n", missing)
	} else {
		fmt.Printf("✅ [DEBUG] すべての環境変数が設定されています\n")
	}
	// ===== デバッグコード終了 =====

	// OAuth1の設定（後々必ず使用するので先にここで行っておく）
	config := oauth1.NewConfig(xApiKey, xPostSecret)
	token := oauth1.NewToken(xAccessToken, xAccessTokenSecret)
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
		fmt.Printf("🔍 [DEBUG] X APIリクエストエラー: %v\n", err)
		return fmt.Errorf("failed to post to x: %w", err)
	}
	defer res.Body.Close()

	// ===== 一時的なデバッグコード（削除予定） =====
	fmt.Printf("🔍 [DEBUG] X APIレスポンス: Status=%s\n", res.Status)
	if res.StatusCode != http.StatusCreated {
		// エラーレスポンスの本文を読み取る
		resBody, readErr := io.ReadAll(res.Body)
		if readErr == nil {
			fmt.Printf("🔍 [DEBUG] エラーレスポンス本文: %s\n", string(resBody))
			return fmt.Errorf("failed to post to x: status %s, body: %s", res.Status, string(resBody))
		}
		fmt.Printf("🔍 [DEBUG] エラーレスポンス本文の読み取りに失敗: %v\n", readErr)
		return fmt.Errorf("failed to post to x: status %s", res.Status)
	}
	fmt.Printf("✅ [DEBUG] X APIへの投稿成功\n")
	// ===== デバッグコード終了 =====

	return nil
}

// ===== 一時的なデバッグ用ヘルパー関数（削除予定） =====
func maskValue(value string) string {
	if value == "" {
		return "(未設定)"
	}
	if len(value) <= 8 {
		return "***"
	}
	// 最初の4文字と最後の4文字を表示
	return value[:4] + "..." + value[len(value)-4:]
}

// ===== デバッグコード終了 =====
