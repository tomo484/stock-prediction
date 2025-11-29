package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"stock-prediction/backend/services"
)

func main() {
	// プロジェクトルートの.envファイルを読み込む
	// backendディレクトリから実行するので ../.env
	envPath := filepath.Join("../../", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("⚠️  警告: .envファイルの読み込みに失敗しました: %v", err)
		log.Println("環境変数から直接読み取りを試みます...")
	} else {
		log.Println("✅ .envファイルを読み込みました")
	}

	// 環境変数からAPIキーを取得
	apiKey := os.Getenv("ALPHA_VANTAGE_API_KEY")
	if apiKey == "" {
		log.Fatal("❌ ALPHA_VANTAGE_API_KEY が設定されていません。\n" +
			"   .envファイルに ALPHA_VANTAGE_API_KEY=your-key を追加するか、\n" +
			"   export ALPHA_VANTAGE_API_KEY=your-key を実行してください。")
	}

	fmt.Println("🔍 Alpha Vantage APIからデータを取得中...")
	fmt.Println("==========================================")

	// Alpha Vantage APIを呼び出し
	data, err := services.FetchAlphaVantageData(apiKey)
	if err != nil {
		log.Fatalf("❌ エラー: %v", err)
	}

	fmt.Println("✅ データ取得成功！")
	fmt.Println("==========================================")

	// 結果を見やすく出力
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("JSON変換エラー: %v", err)
	}

	fmt.Println("\n📊 取得したデータ:")
	fmt.Println(string(jsonData))

	// Top Gainersの件数も表示
	fmt.Printf("\n📈 Top Gainers: %d件\n", len(data.TopGainers))
	fmt.Printf("📉 Top Losers: %d件\n", len(data.TopLosers))
	fmt.Printf("📊 Most Actively Traded: %d件\n", len(data.MostActivelyTraded))

	// サンプルとして最初のTop Gainerを表示
	if len(data.TopGainers) > 0 {
		fmt.Println("\n🏆 Top Gainer #1:")
		fmt.Printf("  Ticker: %s\n", data.TopGainers[0].Ticker)
		fmt.Printf("  Price: $%s\n", data.TopGainers[0].Price)
		fmt.Printf("  Change: %s (%s)\n", data.TopGainers[0].ChangeAmount, data.TopGainers[0].ChangePercentage)
	}
}

