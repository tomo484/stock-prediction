package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"stock-prediction/backend/db"
	"stock-prediction/backend/repositories"
	"stock-prediction/backend/services"
)

func main() {
	// プロジェクトルートの.envファイルを読み込む
	// backend/cmd/test_fmpディレクトリから実行するので ../../../.env
	envPath := filepath.Join("../../../", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("⚠️  警告: .envファイルの読み込みに失敗しました: %v", err)
		log.Println("環境変数から直接読み取りを試みます...")
	} else {
		log.Println("✅ .envファイルを読み込みました")
	}

	// 環境変数からAPIキーを取得
	apiKey := os.Getenv("FMP_API_KEY")
	if apiKey == "" {
		log.Fatal("❌ FMP_API_KEY が設定されていません。\n" +
			"   .envファイルに FMP_API_KEY=your-key を追加するか、\n" +
			"   export FMP_API_KEY=your-key を実行してください。")
	}

	// tickerを取得（引数または環境変数）
	ticker := os.Getenv("TICKER")
	if len(os.Args) > 1 {
		ticker = os.Args[1]
	}
	if ticker == "" {
		log.Fatal("❌ Tickerが指定されていません。\n" +
			"   使用方法: go run main.go AAPL\n" +
			"   または環境変数: TICKER=AAPL go run main.go")
	}

	fmt.Println("🔍 FMP APIから企業情報を取得中...")
	fmt.Printf("   Ticker: %s\n", ticker)
	fmt.Println("==========================================")

	// FMP APIを呼び出し
	fmpData, err := services.FetchFMPData(ticker, apiKey)
	if err != nil {
		log.Fatalf("❌ エラー: %v", err)
	}

	fmt.Println("✅ データ取得成功！")
	fmt.Println("==========================================")

	// 結果を見やすく出力
	jsonData, err := json.MarshalIndent(fmpData, "", "  ")
	if err != nil {
		log.Fatalf("JSON変換エラー: %v", err)
	}

	fmt.Println("\n📊 取得した企業情報:")
	fmt.Println(string(jsonData))

	// 主要な情報を個別に表示
	fmt.Println("\n📈 主要情報:")
	fmt.Printf("  企業名: %s\n", fmpData.CompanyName)
	fmt.Printf("  シンボル: %s\n", fmpData.Symbol)
	fmt.Printf("  セクター: %s\n", fmpData.Sector)
	fmt.Printf("  業界: %s\n", fmpData.Industry)
	fmt.Printf("  現在価格: $%.2f\n", fmpData.Price)
	fmt.Printf("  時価総額: $%.0f\n", fmpData.MarketCap)
	fmt.Printf("  従業員数: %s\n", fmpData.FullTimeEmployees)
	fmt.Printf("  CEO: %s\n", fmpData.CEO)
	fmt.Printf("  ウェブサイト: %s\n", fmpData.Website)
	fmt.Printf("  取引所: %s\n", fmpData.Exchange)

	// DBへの保存テスト（オプション）
	if len(os.Args) > 2 && os.Args[2] == "--save" {
		fmt.Println("\n💾 データベースへの保存をテスト中...")
		
		// データベース接続
		dbConn := db.NewDB()
		defer db.CloseDB(dbConn)

		// Repository初期化
		repo := repositories.NewStockRepository(dbConn)

		// Stockが存在するか確認
		stock, err := repo.FindStockByTicker(ticker)
		if err != nil {
			fmt.Printf("⚠️  Stock %s がDBに存在しません。先にSyncDataを実行してください。\n", ticker)
		} else {
			fmt.Printf("✅ Stock %s が見つかりました (ID: %d)\n", ticker, stock.ID)
			
			// FMPデータをDBに保存
			if err := services.SaveFMPDatatoDB(fmpData, repo); err != nil {
				log.Printf("❌ DB保存エラー: %v", err)
			} else {
				fmt.Println("✅ DBへの保存が完了しました")
			}
		}
	}

	fmt.Println("\n✅ テスト完了！")
}


