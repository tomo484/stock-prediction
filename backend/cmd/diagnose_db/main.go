package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// .envファイルを読み込む
	envPath := filepath.Join("..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("⚠️  .envファイルが見つかりません: %v\n", err)
	} else {
		fmt.Println("✅ .envファイルを読み込みました")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	fmt.Println("\n📊 データベース設定:")
	fmt.Printf("  Host: %s\n", host)
	fmt.Printf("  User: %s\n", user)
	fmt.Printf("  Password: %s (長さ: %d文字)\n", maskPassword(password), len(password))
	fmt.Printf("  Database: %s\n", dbname)
	fmt.Printf("  Port: %s\n", port)
	fmt.Println("==========================================")

	// まず、postgresデータベースに接続を試みる（デフォルトDB）
	fmt.Println("\n🔍 ステップ1: デフォルトDB (postgres) への接続を試行...")
	postgresConnStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		host, port, user, password,
	)
	
	db, err := sql.Open("postgres", postgresConnStr)
	if err != nil {
		log.Printf("❌ 接続文字列の作成に失敗: %v\n", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Printf("❌ postgresデータベースへの接続に失敗しました\n")
		fmt.Printf("   エラー: %v\n\n", err)
		fmt.Println("📝 考えられる原因:")
		fmt.Println("   1. PostgreSQLが起動していない")
		fmt.Println("   2. パスワードが間違っている")
		fmt.Println("   3. ユーザー名が間違っている")
		fmt.Println("   4. ホスト/ポートが間違っている")
		fmt.Println("\n💡 解決方法:")
		fmt.Println("   sudo service postgresql status  # PostgreSQLの状態確認")
		fmt.Println("   sudo service postgresql start   # PostgreSQLを起動")
		fmt.Println("   sudo -u postgres psql           # パスワードなしでログイン")
		return
	}
	
	fmt.Println("✅ postgresデータベースへの接続成功！")

	// データベース一覧を取得
	fmt.Println("\n📋 ステップ2: データベース一覧を確認...")
	rows, err := db.Query("SELECT datname FROM pg_database WHERE datistemplate = false")
	if err != nil {
		log.Printf("❌ データベース一覧の取得に失敗: %v\n", err)
		return
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			continue
		}
		databases = append(databases, dbName)
	}

	fmt.Println("  利用可能なデータベース:")
	targetExists := false
	for _, dbName := range databases {
		if dbName == dbname {
			fmt.Printf("  ✅ %s (ターゲットDB)\n", dbName)
			targetExists = true
		} else {
			fmt.Printf("  - %s\n", dbName)
		}
	}

	// ターゲットDBが存在しない場合
	if !targetExists {
		fmt.Printf("\n⚠️  ターゲットデータベース '%s' が見つかりません\n", dbname)
		fmt.Println("\n💡 データベースを作成するには:")
		fmt.Println("   sudo -u postgres psql")
		fmt.Printf("   CREATE DATABASE %s;\n", dbname)
		fmt.Println("   \\q")
		
		// 自動作成を提案
		fmt.Println("\n🔧 データベースを自動作成しますか？")
		createQuery := fmt.Sprintf("CREATE DATABASE %s", dbname)
		_, err := db.Exec(createQuery)
		if err != nil {
			fmt.Printf("❌ データベースの作成に失敗: %v\n", err)
			return
		}
		fmt.Printf("✅ データベース '%s' を作成しました！\n", dbname)
		targetExists = true
	}

	// ターゲットDBへの接続を試みる
	if targetExists {
		fmt.Printf("\n🔍 ステップ3: ターゲットDB (%s) への接続を試行...\n", dbname)
		targetConnStr := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname,
		)
		
		targetDB, err := sql.Open("postgres", targetConnStr)
		if err != nil {
			log.Printf("❌ 接続に失敗: %v\n", err)
			return
		}
		defer targetDB.Close()

		err = targetDB.Ping()
		if err != nil {
			fmt.Printf("❌ ターゲットデータベースへの接続に失敗: %v\n", err)
			return
		}
		
		fmt.Printf("✅ ターゲットデータベース '%s' への接続成功！\n", dbname)
		
		// テーブル一覧を取得
		fmt.Println("\n📊 ステップ4: テーブル一覧を確認...")
		tableRows, err := targetDB.Query("SELECT tablename FROM pg_tables WHERE schemaname = 'public'")
		if err != nil {
			log.Printf("❌ テーブル一覧の取得に失敗: %v\n", err)
			return
		}
		defer tableRows.Close()

		var tables []string
		for tableRows.Next() {
			var tableName string
			if err := tableRows.Scan(&tableName); err != nil {
				continue
			}
			tables = append(tables, tableName)
		}

		if len(tables) == 0 {
			fmt.Println("  ⚠️  テーブルが見つかりません")
			fmt.Println("\n💡 マイグレーションを実行してテーブルを作成してください:")
			fmt.Println("   go run main.go")
		} else {
			fmt.Println("  存在するテーブル:")
			for _, table := range tables {
				fmt.Printf("  - %s\n", table)
			}
		}
	}

	fmt.Println("\n✅ 診断完了！")
}

func maskPassword(password string) string {
	if password == "" {
		return "(空)"
	}
	if len(password) <= 4 {
		return "****"
	}
	return password[:2] + "****" + password[len(password)-2:]
}


