package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"stock-prediction/backend/db"
	"stock-prediction/backend/models"
	"stock-prediction/backend/repositories"
	xpost "stock-prediction/backend/services/x_post"

	"github.com/joho/godotenv"
)

func main() {
	// プロジェクトルートの.envファイルを読み込む
	envPath := filepath.Join("../../../", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("⚠️  警告: .envファイルの読み込みに失敗しました: %v", err)
		log.Println("環境変数から直接読み取りを試みます...")
	} else {
		log.Println("✅ .envファイルを読み込みました")
	}

	// テストタイプを取得（引数から）
	testType := "all"
	if len(os.Args) > 1 {
		testType = os.Args[1]
	}

	// 日付を取得（オプション）
	date := ""
	if len(os.Args) > 2 {
		date = os.Args[2]
	}

	fmt.Println("🔍 X自動投稿機能のテスト（Dry-Runモード）")
	fmt.Println("==========================================")

	// データベース接続
	fmt.Println("\n📊 データベースに接続中...")
	dbConn := db.NewDB()
	defer db.CloseDB(dbConn)

	// Repository初期化
	repo := repositories.NewStockRepository(dbConn)

	// 最新の日付を取得（指定されていない場合）
	if date == "" {
		var latestDate string
		result := dbConn.Model(&models.DailyRanking{}).
			Select("MAX(date)").
			Scan(&latestDate)

		if result.Error != nil || latestDate == "" {
			log.Fatal("❌ データベースから日付を取得できませんでした。先にSyncDataを実行してください。")
		}
		date = latestDate
		fmt.Printf("✅ 最新日付を使用: %s\n", date)
	} else {
		fmt.Printf("✅ 指定日付を使用: %s\n", date)
	}

	// テストタイプに応じて実行
	switch testType {
	case "ranking":
		testRankingPost(date, repo)
	case "analysis":
		testAnalysisPost(date, repo)
	case "all":
		testRankingPost(date, repo)
		fmt.Println("\n" + strings.Repeat("=", 50) + "\n")
		testAnalysisPost(date, repo)
	default:
		log.Fatalf("❌ 無効なテストタイプ: %s\n   使用方法: go run main.go [ranking|analysis|all] [date]", testType)
	}

	fmt.Println("\n✅ テスト完了！")
	fmt.Println("⚠️  注意: これはDry-Runモードです。実際のXへの投稿は行われていません。")
}

func testRankingPost(date string, repo repositories.IStockRepository) {
	fmt.Println("\n📝 ランキング投稿のテスト")
	fmt.Println("----------------------------------------")

	// ランキングデータを取得
	rankings, err := repo.FindDailyRanking(date)
	if err != nil {
		log.Fatalf("❌ ランキングデータの取得に失敗: %v", err)
	}

	if len(*rankings) == 0 {
		log.Fatal("❌ ランキングデータが見つかりませんでした")
	}

	// テンプレート生成
	text := xpost.BuildRankingPost(date, *rankings)

	// 結果を表示
	fmt.Println("\n📤 生成された投稿内容:")
	fmt.Println("----------------------------------------")
	fmt.Println(text)
	fmt.Println("----------------------------------------")

	// 文字数チェック
	charCount := utf8.RuneCountInString(text)
	fmt.Printf("\n📊 文字数: %d / 280文字\n", charCount)
	if charCount > 280 {
		fmt.Printf("⚠️  警告: 文字数制限（280文字）を超過しています！\n")
	} else {
		fmt.Printf("✅ 文字数制限内です（残り: %d文字）\n", 280-charCount)
	}
}

func testAnalysisPost(date string, repo repositories.IStockRepository) {
	fmt.Println("\n📝 個別分析投稿のテスト")
	fmt.Println("----------------------------------------")

	// 1位から5位までテスト
	for rank := 1; rank <= 5; rank++ {
		fmt.Printf("\n🏆 Rank %d の投稿テスト:\n", rank)
		fmt.Println("----------------------------------------")

		// ランキングデータを取得
		ranking, err := repo.FindDailyRankingByDateAndRank(date, rank, "Top Gainers")
		if err != nil {
			fmt.Printf("⚠️  Rank %d のデータが見つかりませんでした: %v\n", rank, err)
			continue
		}

		// AiAnalysisが空の場合はスキップ
		if ranking.AiAnalysis == "" {
			fmt.Printf("⚠️  Rank %d のAiAnalysisが空です。スキップします。\n", rank)
			continue
		}

		// テンプレート生成
		text := xpost.BuildAnalysisPost(*ranking)

		// 結果を表示
		fmt.Println("\n📤 生成された投稿内容:")
		fmt.Println("----------------------------------------")
		fmt.Println(text)
		fmt.Println("----------------------------------------")

		// 文字数チェック
		charCount := utf8.RuneCountInString(text)
		fmt.Printf("\n📊 文字数: %d / 280文字\n", charCount)
		if charCount > 280 {
			fmt.Printf("⚠️  警告: 文字数制限（280文字）を超過しています！\n")
		} else {
			fmt.Printf("✅ 文字数制限内です（残り: %d文字）\n", 280-charCount)
		}

		// 元のAiAnalysisの文字数も表示
		originalLen := utf8.RuneCountInString(ranking.AiAnalysis)
		if originalLen > 280 {
			fmt.Printf("📝 元のAiAnalysis: %d文字（切り詰められました）\n", originalLen)
		}
	}
}
