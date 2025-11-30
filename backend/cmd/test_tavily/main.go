package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"stock-prediction/backend/services/news"

	"github.com/joho/godotenv"
)

func main() {
	// コマンドライン引数のチェック
	if len(os.Args) < 2 {
		log.Fatal("❌ エラー: tickerが指定されていません。\n" +
			"   使用方法: go run main.go <TICKER>\n" +
			"   例: go run main.go AAPL")
	}

	ticker := os.Args[1]
	fmt.Printf("🔍 Ticker: %s のニュースを検索中...\n", ticker)
	fmt.Println("==========================================")

	// プロジェクトルートの.envファイルを読み込む
	// backend/cmd/test_tavilyディレクトリから実行するので ../../../.env
	envPath := filepath.Join("../../../", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("⚠️  警告: .envファイルの読み込みに失敗しました: %v", err)
		log.Println("環境変数から直接読み取りを試みます...")
	} else {
		log.Println("✅ .envファイルを読み込みました")
	}

	// 環境変数からAPIキーを取得
	apiKey := os.Getenv("TAVILY_API_KEY")
	if apiKey == "" {
		log.Fatal("❌ TAVILY_API_KEY が設定されていません。\n" +
			"   .envファイルに TAVILY_API_KEY=your-key を追加するか、\n" +
			"   export TAVILY_API_KEY=your-key を実行してください。")
	}

	// Tavily Search APIを呼び出し
	headlines, err := news.SearchStockNews(ticker, apiKey)
	if err != nil {
		log.Fatalf("❌ エラー: %v", err)
	}

	fmt.Println("✅ 検索成功！")
	fmt.Println("==========================================")

	// 結果を表示
	if len(headlines) == 0 {
		fmt.Printf("\n⚠️  %s に関するニュースが見つかりませんでした。\n", ticker)
		return
	}

	fmt.Printf("\n📰 取得したニュース: %d件\n\n", len(headlines))
	for i, headline := range headlines {
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("📄 ニュース #%d\n", i+1)
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Println(headline)
		fmt.Println()
	}
}
