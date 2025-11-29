package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"stock-prediction/backend/models"
)

func main() {
	// .envファイルを読み込む
	envPath := filepath.Join("..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("⚠️  警告: .envファイルの読み込みに失敗しました: %v", err)
		log.Println("環境変数から直接読み取りを試みます...")
	} else {
		log.Println("✅ .envファイルを読み込みました")
	}

	// データベース接続情報を表示（パスワードは隠す）
	fmt.Println("\n📊 データベース接続情報:")
	fmt.Printf("  Host: %s\n", os.Getenv("DB_HOST"))
	fmt.Printf("  User: %s\n", os.Getenv("DB_USER"))
	fmt.Printf("  Password: %s\n", hidePassword(os.Getenv("DB_PASSWORD")))
	fmt.Printf("  Database: %s\n", os.Getenv("DB_NAME"))
	fmt.Printf("  Port: %s\n", os.Getenv("DB_PORT"))
	fmt.Println("==========================================")

	// PostgreSQL接続URLを構築
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// データベースに接続
	fmt.Println("\n🔌 データベースに接続中...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ データベース接続エラー: %v", err)
	}
	fmt.Println("✅ データベースに接続しました")

	// マイグレーションを実行
	fmt.Println("\n🔧 マイグレーションを実行中...")
	if err := db.AutoMigrate(&models.Stock{}, &models.DailyRanking{}); err != nil {
		log.Fatalf("❌ マイグレーションエラー: %v", err)
	}
	fmt.Println("✅ マイグレーション完了")

	// テーブル一覧を取得
	fmt.Println("\n📋 テーブル一覧:")
	var tables []string
	db.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables)
	if len(tables) == 0 {
		fmt.Println("  ⚠️  テーブルが見つかりませんでした")
	} else {
		for _, table := range tables {
			fmt.Printf("  - %s\n", table)
		}
	}

	// Stockテーブルのレコード数を確認
	fmt.Println("\n📊 テーブルのレコード数:")
	var stockCount int64
	db.Model(&models.Stock{}).Count(&stockCount)
	fmt.Printf("  stocks: %d件\n", stockCount)

	var rankingCount int64
	db.Model(&models.DailyRanking{}).Count(&rankingCount)
	fmt.Printf("  daily_rankings: %d件\n", rankingCount)

	// Stockテーブルのデータをサンプル表示（最初の5件）
	fmt.Println("\n📦 Stockテーブルのサンプルデータ (最初の5件):")
	var stocks []models.Stock
	db.Limit(5).Find(&stocks)
	if len(stocks) == 0 {
		fmt.Println("  データがありません")
	} else {
		for i, stock := range stocks {
			fmt.Printf("  %d. Ticker: %s, Name: %s, Sector: %s\n",
				i+1, stock.Ticker, stock.Name, stock.Sector)
		}
	}

	// DailyRankingテーブルのデータをサンプル表示（最初の5件）
	fmt.Println("\n📈 DailyRankingテーブルのサンプルデータ (最初の5件):")
	var rankings []models.DailyRanking
	db.Limit(5).Find(&rankings)
	if len(rankings) == 0 {
		fmt.Println("  データがありません")
	} else {
		for i, ranking := range rankings {
			fmt.Printf("  %d. Date: %s, Rank: %d, StockID: %d, Price: $%.2f, ChangeRate: %.2f%%\n",
				i+1, ranking.Date, ranking.Rank, ranking.StockID, ranking.Price, ranking.ChangeRate)
		}
	}

	fmt.Println("\n✅ データベース確認完了！")
}

// hidePassword はパスワードを部分的に隠す
func hidePassword(password string) string {
	if password == "" {
		return "(空)"
	}
	if len(password) <= 4 {
		return "****"
	}
	return password[:2] + "****" + password[len(password)-2:]
}

