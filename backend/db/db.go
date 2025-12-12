package db

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDB はPostgreSQLデータベースへの接続を確立し、GORMのDBインスタンスを返す
func NewDB() *gorm.DB {
	// プロジェクトルートの.envファイルを読み込む
	// backendディレクトリから実行する場合、../.env の位置
	envPath := filepath.Join("..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		// カレントディレクトリの.envも試す（直接プロジェクトルートから実行する場合）
		if err2 := godotenv.Load(); err2 != nil {
			log.Println("Warning: .env file not found, using environment variables")
		}
	}

	// PostgreSQL接続URLを取得
	// supabaseDB_URL環境変数から接続情報を取得
	dsn := os.Getenv("supabaseDB_URL")
	if dsn == "" {
		log.Fatalln("supabaseDB_URL environment variable is not set")
	}

	// GORMでPostgreSQLに接続
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Failed to connect database:", err)
	}

	log.Println("Connected to database")
	return db
}

// CloseDB はデータベース接続を閉じる
func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalln("Failed to get database instance:", err)
	}

	if err := sqlDB.Close(); err != nil {
		log.Fatalln("Failed to close database:", err)
	}

	log.Println("Database connection closed")
}
