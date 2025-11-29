(## システムを作る上での基本仕様・技術スタック
・Go(echo) , TypeScript (next.js),gorm でシステムを作る
・システムの中身は、前日の米国株の中で最も上昇率が高かった銘柄5選の上昇率の理由をAIに調査させるという形を想定している
・日付は全て米国基準（米国時間）で扱う
・株の情報を取得する部分は、Alpha VantageのAPIを使おうと思ってる
・AIのAPIエンドポイントは普通にOpen AIのやつでいいかな
・バックエンドのアーキテクチャはcontroller層, repositorya層, service層の3層のクリーンアーキテクチャで行う

## DB設計
-PostgreSQLとgormを使用する
-データベースの構造（データテーブル構造）配下の通り

### Stock（株価銘柄基本情報）
type Stock struct {
    gorm.Model
    Ticker      string `gorm:"uniqueIndex;not null" json:"ticker"` // 銘柄コード (例: NVDA)
    Name        string `json:"name"`                               // 企業名
    Sector      string `json:"sector"`                             // セクター (例: Technology) ※将来の分析用
    Industry    string `json:"industry"`                           // 業界 (例: Semiconductors)
    
    // リレーション設定: 1つの銘柄は、多数のランキング履歴を持つ
    Rankings    []DailyRanking `gorm:"foreignKey:StockID" json:"rankings,omitempty"` 
}

### Ranking（ランキング）
type DailyRanking struct {
    gorm.Model
    StockID     uint    `json:"stock_id"`      // Foreign Key (Stockテーブルへの紐付け)
    Date        string  `gorm:"index" json:"date"` // 日付 (例: "2025-01-01") ※米国基準（米国時間）
    Rank        int     `json:"rank"`          // その日の順位
    ChangeRate  float64 `json:"change_rate"`   // 上昇率
    Price       float64 `json:"price"`         // その時の株価
    
    // AI関連
    NewsSummary string  `gorm:"type:text" json:"news_summary"` // AIに読ませたニュースの要約（念のため保存）
    AiAnalysis  string  `gorm:"type:text" json:"ai_analysis"`  // AIが出した上昇理由
    
    // リレーション設定: ランキングは1つの銘柄に紐づく
    Stock       Stock   `json:"stock,omitempty"` 
}

## バックエンド側（APIエンドポイント）

### フロントエンド表示用API
HTTPメソッド,パス (URL),役割,処理内容
GET,/api/stocks/latest,最新ランキング取得,DBから最新日付のトップ5銘柄とAI解説を取得して返す。（とにかく一番最新のやつを拾ってくる = DailyRankingのIDカラムの最新のIDを拾ってくればいい）
GET,/api/stocks/date,日付指定ランキング取得,クエリパラメータで日付を指定してその日のランキングを取得する。例: /api/stocks/date?date=2025-01-15（日付は米国基準、フォーマット: YYYY-MM-DD）
GET,/api/stocks/:ticker,過去の銘柄ログ取得 (任意),ティッカーコードを指定して過去のランキングを取得する（カレンダー機能用など）。返却データはそのStockの過去のDailyRankingの配列。

### データ取得・管理用
HTTPメソッド,パス (URL),役割,処理内容
POST,/api/admin/sync,データ収集トリガー,重要: これが「Alpha Vantage取得 + AI分析」の実行ボタンになります。Alpha Vantage API叩く → 上昇率計算 → Top5抽出 → AIに投げる → DBに保存。（取得したTickerの中で新しく登場したものは新しくStockテーブルのﾚｺｰﾄﾞとして登録する。既にﾚｺｰﾄﾞがあるTicerの企業が登場した場合は既存銘柄の更新（Name, Sector, Industry）は行う必要はない。
PUT,/api/stocks/rankings/:id,AI回答の修正,AIの解説がおかしい時に手動で書き換える用。（変更を加えるのはDailyRankingのAiAnalysisのみ。他には手を加えない。:idはDailyRankingのID）
DELETE,/api/stocks/rankings/:id,データ削除,バグで変なデータが入った時などに消す用。（:idはDailyRankingのID）

### フロントエンド側表示用APIのcontroler層, repository層, service層の引数・返り値

#### /api/stocks/date（GET）

##### Controller層
-引数：c echo.Context（クエリパラメータ `date` から日付を取得。例: `c.QueryParam("date")`）
-返り値：error

##### Service層
-引数：Date string（フォーマット: YYYY-MM-DD、米国基準）
-返り値：*[]models.DailyRanking, error

#### Repository層
-引数：Date string（フォーマット: YYYY-MM-DD、米国基準）
-返り値：*[]models.DailyRanking, error

#### /api/stocks/latest（GET）

##### Controller層
-引数：c echo.Context
-返り値：error

##### Service層
-引数：()
-返り値：*[]models.DailyRanking, error

#### Repository層
-引数：()
-返り値：*[]models.DailyRanking, error

#### /api/stocks/:ticker（GET）

##### Controller層
-引数：c echo.Context
-返り値：error

##### Service層
-引数：Ticker string
-返り値：*[]models.DailyRanking, error

#### Repository層
-引数：Ticker string 
-返り値：*[]models.DailyRanking, error

#### /api/admin/sync

##### Controller層
-引数：c echo.Context
-返り値：error

##### Service層
-引数：dto.AlphaVantageRequest（Alpha-VantageのAPIにアクセスするときに必要な引数を定義したもの）
-返り値：*[]models.DailyRanking, error

#### Repository層
-引数：[]models.DailyRanking
-返り値：*[]models.DailyRanking, error

#### api/stocks/rankings/:id

##### Controller層
-引数：c echo.Context
-返り値：error

##### Service層
-引数：ID uint, AiAnalysis string
-返り値：*models.DailyRanking, error

#### Repository層
-引数：ID uint, AiAnalysis string
-返り値：*models.DailyRanking, error

#### /api/stocks/rankings/:id

##### Controller層
-引数：c echo.Context
-返り値：error

##### Service層
-引数：ID uint 
-返り値：error

#### Repository層
-引数：ID uint
-返り値：error


## フロントエンド側状態管理（ReactQueryカスタムフック）

### 表示・取得用(GET)
-画面にデータを表示するためのフック

useLatestStocks (GET /api/stocks/latest)
-最新のランキングを表示する用。
querykey:["ranking", "latest"]
引数：()

useDateStocks (GET /api/stocks/date)
-dateを指定して、その日のランキングを表示する用。API呼び出しは `/api/stocks/date?date=YYYY-MM-DD` の)形式でクエリパラメータを使用する。
querykey:["ranking", date]
引数：date: string (フォーマット: YYYY-MM-DD、米国基準)

useStockHistory (GET /api/stocks/:ticker)
-過去データを見る用（引数で Ticker string を受け取る）。
querykey:["stocks", ticker]
引数：ticker: string

### 操作・更新用
-ボタンを押してアクションをおこす為のフックです。

useSyncStocks (POST /api/admin/sync)
-「Alpha Vantage取得＆AI分析」ボタン用。
invalidateQueries:["ranking"]
引数：()

useUpdateStockAi (PUT /api/stocks/rankings/:id)
-AIのコメント修正用。（引数で DailyRankingのID を受け取る）
invalidateQueries:["ranking", date]->["stocks", ticker]
引数：id: number, Aianalysis: string

useDeleteStock (DELETE /api/stocks/rankings/:id)
-データ削除用。（引数で DailyRankingのID を受け取る）
invalidateQueries:["ranking", date]->["stocks", ticker]
引数：id: number

## 株価取得API（米国株）
-以下のAPIを株価取得APIとして使用する
・Alpha Vantage (Top Gainers API)
-API名: Alpha Vantage (TOP_GAINERS_LOSERS endpoint)
-費用: 無料（1日25リクエストまで。今回の「1日1回ランキング取得」なら余裕です）
-取得できるデータ: 値上がり上位20銘柄の「ティッカー」「価格」「上昇率」「取引量」

## 画面構成（UI部分）について
3ページ構成
・/page.tsx -HPみたいな感じのおしゃれなページ（ここでは上昇ランキングの銘柄を表示するだけで十分）
・/dashboard -あなたが提案してくれたメインダッシュボードページ
・/stocks/[ticker] -あなたが提案してくれた銘柄詳細ページ

### ホームページ(/page.tsx)
① 機能（Function）
ヒーローセクション: サイトのコンセプト（「AIが急騰株を即座に分析」）を伝えるキャッチコピー。
本日のTop 5（簡易版）: AI分析文は載せるが、出来れば途中で研ぎらせれるような魅せ方がよく、「順位・銘柄・上昇率」だけをシンプルに見せる。
CTA (Call To Action): 「AIによる分析を読む（Dashboardへ）」ボタン。

②ワイヤーフレーム
+----------------------------------------------------+
| [Logo]                                 [Dashboard] |
+----------------------------------------------------+
|                                                    |
|      今日の米国株市場、                          |
|      なぜその株が上がったのか？                  |
|      AIが一瞬で解明します。                      |
|                                                    |
|         [  AI分析を見る (Button)  ]              |
|                                                    |
+----------------------------------------------------+
|  Today's Top Gainers (2025-11-24)                  |
|                                                    |
|  1. 🥇 NVIDIA (NVDA) .......... +5.4%             |
|  2. 🥈 Apple (AAPL) .......... +4.2%             |
|  3. 🥉 Microsoft (MSFT) .......... +3.8%             |
|  4. 4  Tesla (TSLA) ...... +2.1%             |
|  5. 5  Amazon (AMZN) .... +1.5%             |
|                                                    |
+----------------------------------------------------+
| (Footer: 免責事項...)                              |
+----------------------------------------------------+
デザインのポイント: 背景にうっすらとチャートの曲線をあしらったり、数字（+5.4%）を緑色や赤色で光らせるなどして、テック感を出します。

### メインダッシュボード（/dashboard/page.tsx）
①機能（Function）
URL形式: `/dashboard` (最新データ) または `/dashboard?date=2025-01-15` (指定日付)。日付は米国基準、フォーマット: YYYY-MM-DD
日付ナビゲーション: デフォルトは「今日（最新データ）」。カレンダーまたは「＜ 前日 翌日 ＞」ボタンで過去のランキングに切り替え可能。日付変更時はURLのクエリパラメータを更新する。
詳細カードリスト: Top 5銘柄それぞれの「AI分析結果」を表示。
ドリルダウン: 銘柄名をクリックすると、その銘柄の詳細ページへ遷移。

②ワイヤーフレーム
+----------------------------------------------------+
| [Logo]                                             |
+----------------------------------------------------+
|  [<]    2025年11月24日 (月)    [>]  [📅]         |
+----------------------------------------------------+
|                                                    |
| +------------------------------------------------+ |
| | 🥇 1位: NVIDIA Corporation (NVDA)          +5.4% 📈 | |
| |------------------------------------------------| |
| | 【AI分析サマリー】                             | |
| |  AIチップ需要の急拡大とデータセンター向けGPU売上増加が...      | |
| |  特に生成AI市場の成長見通しが上方修正されたニュースが...     | |
| |  (ニュースソース: Bloomberg, Reuters)           | |
| |                                                | |
| |              [ 過去の履歴を見る > ]            | |
| +------------------------------------------------+ |
|                                                    |
| +------------------------------------------------+ |
| | 🥈 2位: Apple Inc. (AAPL)        +4.2% 📈 | |
| |------------------------------------------------| |
| | 【AI分析サマリー】                             | |
| |  .....                                         | |
| +------------------------------------------------+ |
|   (以下、3位〜5位まで続く)                         |
+----------------------------------------------------+
UIコンポーネント: アコーディオン（開閉式）にはせず、最初から展開して読ませる形が良いでしょう。スクロールで流し読みさせます。

③銘柄詳細ページ
① 機能（Function）
銘柄基本情報: 企業名、ティッカー、セクター（DBにあれば）。
ランクイン履歴 (Timeline): このシステムで「Top 5に入った日」のリスト。
いつ、何位で、どのくらい上がって、その時AIは何と言ったか。
株価の特性（決算で上がるタイプか、材料で飛ぶタイプか）が見えてきます。

②ワイヤーフレーム
+----------------------------------------------------+
| [Logo]  < Dashboardに戻る                          |
+----------------------------------------------------+
|  [ 🏭 Technologyセクター ]                             |
|  NVIDIA Corporation (NVDA)                               |
|  NVIDIA CORP.                                |
+----------------------------------------------------+
|  History (ランクイン履歴)                          |
|                                                    |
|  ● 2025-11-24 (1位 / +5.4%)                        |
|    | 理由: AIチップ需要の急拡大とデータセンター向けGPU売上増加...       |
|    +-------------------------------------------    |
|                                                    |
|  ● 2025-08-10 (3位 / +3.1%)                        |
|    | 理由: 第1四半期決算がコンセンサス予想を上回る...    |
|    +-------------------------------------------    |
|                                                    |
|  ● 2025-05-22 (5位 / +2.0%)                        |
|    | 理由: 新製品発表とパートナーシップ拡大のニュース...         |
|    +-------------------------------------------    |
+----------------------------------------------------+

## ディレクトリ形式
-ディレクトリ形式は以下のような形式で行う
my-stock-app/
├── backend/             (Go + Echo)
│   ├── go.mod
│   ├── main.go          (エントリーポイント)
│   ├── db/
│   │   └── db.go        (GORM接続設定:postgresql)
│   ├── models/
│   │   └── stock.go     (Stock, DailyRankingのStruct定義)
│   ├── controllers/
│   │   └── stock_controller.go     (APIのエンドポイント処理: GetStock, SyncDataなど, Controller層。admin関連のエンドポイントも含む)
│   ├── services/
│   │   ├── stock_service.go        (service層)
│   │   ├── openai.go               (ChatGPT APIへの問い合わせロジック: API呼び出しのみ)
│   │   ├── alpha.go                (Alpha-Vantageのranking APIへの問い合わせロジック)
│   │   └── AIanalysis.go           (システムプロンプト等のAI対応について書き込むファイル: システムプロンプトと分析ロジック)
│   ├── repositories/
│   │   └── stock_repository.go     (Repository層)
│   ├── router/
│   │   └── router.go               (Echoルーティング設定、ミドルウェア設定)
│   └── dto/
│       └── stock_dto.go            (dto定義を書く場所)
|
│
└── frontend/            (Next.js + TypeScript)
    ├── package.json
    ├── next.config.mjs
    └── src/
        ├── app/
        │   ├── page.tsx             (LP: トップページ)
        │   ├── dashboard/
        │   │   └── page.tsx         (メイン: ランキング一覧 & 更新ボタン)
        │   └── stocks/
        │       └── [ticker]/
        │           └── page.tsx     (詳細: 過去履歴)
        ├── components/
        │   ├── StockCard.tsx        (ランキング表示用カード)
        │   ├── SyncButton.tsx       (更新ボタン: ローディング状態管理)
        │   └── ui/                  (shadcn/ui のボタン等が入る)
        ├── hooks/
        │   └── useStocks.ts         (React Queryのカスタムフック集)
        ├── lib/
        │   └── axios.ts             (APIクライアント設定)
        └── types/
            └── stock.ts             (stockの型定義を行う)

## その他議論が必要な事柄に関して
-/dashboard.page.tsxに株価取得APIを走らせるボタンをつける。但しこの仕様は初期のみで使い慣れてきたら1日1回のイベント駆動型にする
-AIの部分に関して、Langchain, Langgraphは開発初期（現時点）では使用せず、機能を拡張する際に使用することとする。基本的にシステムプロンプトだけで代用するが、将来的にはこのNewsのAPIをたたいて、情報を取ってきてみたいな挙動もさせたいのでそうなるとLanggraphは必須になるから拡張性を残した実装をしたいね
-株価取得APIはAlpha Vantageを使用するけど、 alpha.goのファイルで基本的にapiアクセスは行うつもり
-