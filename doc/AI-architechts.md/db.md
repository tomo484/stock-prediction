# 日本株分析システム DB設計

## 概要

Go (Backend) で J-Quants API および Tavily API からデータを収集し、PostgreSQL (Supabase) に保存する。
Python (AI Agent) には gRPC 経由でデータを渡し、LangGraph で分析を行う。

---

## テーブル一覧

| # | テーブル名 | 説明 | データソース |
|---|------------|------|--------------|
| 1 | companies | 銘柄マスタ | J-Quants /listed/info |
| 2 | daily_quotes | 日足株価（過去6ヶ月分） | J-Quants /prices/daily_quotes |
| 3 | financial_statements | 財務諸表（過去5年分） | J-Quants /fins/statements |
| 4 | news_searches | ニュース検索バッチ | Tavily Search API |
| 5 | news_items | ニュース検索結果（個別） | Tavily Search API |
| 6 | analysis_results | AI分析結果 | Python gRPC Response |

---

## 1. companies (銘柄マスタ)

J-Quants `/listed/info` から取得した銘柄情報を保存する。

### APIレスポンス例
```json
{
  "info": [
    {
      "Date": "2022-11-11",
      "Code": "86970",
      "CompanyName": "日本取引所グループ",
      "CompanyNameEnglish": "Japan Exchange Group,Inc.",
      "Sector17Code": "16",
      "Sector17CodeName": "金融（除く銀行）",
      "Sector33Code": "7200",
      "Sector33CodeName": "その他金融業",
      "ScaleCategory": "TOPIX Large70",
      "MarketCode": "0111",
      "MarketCodeName": "プライム"
    }
  ]
}
```

### Goモデル定義
```go
type Company struct {
    ID                  uint   `gorm:"primaryKey"`
    Code                string `gorm:"uniqueIndex;not null"`  // 銘柄コード（例: "86970"）
    CompanyName         string                                // 企業名
    CompanyNameEnglish  string                                // 企業名（英語）
    Sector17Code        string `gorm:"index"`                 // 17業種コード
    Sector17CodeName    string                                // 17業種名
    Sector33Code        string `gorm:"index"`                 // 33業種コード
    Sector33CodeName    string                                // 33業種名
    ScaleCategory       string                                // 規模区分
    MarketCode          string `gorm:"index"`                 // 市場コード
    MarketCodeName      string                                // 市場名
    UpdatedAt           time.Time                             // 最終更新日時
}
```

### 設計メモ
- `Code` を一意制約とする（主キーは別途 `ID` を持つ）
- セクター名・市場名の正規化は行わない（マスタ更新頻度が低いため）

---

## 2. daily_quotes (日足株価)

J-Quants `/prices/daily_quotes` から取得した日足株価データを保存する。
過去6ヶ月分を保持。

### APIレスポンス例
```json
{
  "daily_quotes": [
    {
      "Date": "2023-03-24",
      "Code": "86970",
      "Open": 2047.0,
      "High": 2069.0,
      "Low": 2035.0,
      "Close": 2045.0,
      "Volume": 2202500.0,
      "TurnoverValue": 4507051850.0,
      "AdjustmentFactor": 1.0,
      "AdjustmentOpen": 2047.0,
      "AdjustmentHigh": 2069.0,
      "AdjustmentLow": 2035.0,
      "AdjustmentClose": 2045.0,
      "AdjustmentVolume": 2202500.0
    }
  ]
}
```

### Goモデル定義
```go
type DailyQuote struct {
    ID               uint    `gorm:"primaryKey"`
    Code             string  `gorm:"index:idx_daily_quote_code_date,unique;not null"`
    Date             string  `gorm:"index:idx_daily_quote_code_date,unique;not null"`
    Open             float64
    High             float64
    Low              float64
    Close            float64
    Volume           float64
    TurnoverValue    float64                          // 売買代金
    AdjustmentFactor float64                          // 調整係数
    AdjustmentOpen   float64                          // 調整後始値
    AdjustmentHigh   float64                          // 調整後高値
    AdjustmentLow    float64                          // 調整後安値
    AdjustmentClose  float64                          // 調整後終値
    AdjustmentVolume float64                          // 調整後出来高
}
```

### 設計メモ
- `ID` を主キーとし、`(Code, Date)` にユニーク制約を付ける
- GORMでの扱いやすさと他テーブルからの参照を考慮してサロゲートキー方式を採用

---

## 3. financial_statements (財務諸表)

J-Quants `/fins/statements` から取得した財務諸表データを保存する。
過去5年分（四半期・通期全て）を保持。

### 設計方針: ハイブリッド方式

**理由:**
- 財務諸表は70カラム以上あり、全てをカラム定義すると保守コストが高い
- Python側で `pandas.read_json` で解析する設計のため、生JSONをそのまま渡せる方が効率的
- ただし、フィルタリング（期種別、開示日等）はGo側でも行いたい

**方針:**
- 検索・フィルタに使うフィールドのみカラム化
- 残りは `RawJSON` として保存

### Goモデル定義
```go
type FinancialStatement struct {
    ID                  uint   `gorm:"primaryKey"`
    
    // === 検索・フィルタ用（カラム化） ===
    Code                string `gorm:"index;not null"`           // 銘柄コード
    DisclosureNumber    string `gorm:"uniqueIndex;not null"`     // 開示番号（一意識別子）
    DisclosedDate       string `gorm:"index"`                    // 開示日
    TypeOfDocument      string `gorm:"index"`                    // 書類種別
    TypeOfCurrentPeriod string                                   // 期種別（1Q/2Q/3Q/FY）
    CurrentFiscalYearStartDate string                            // 会計年度開始日
    CurrentFiscalYearEndDate   string `gorm:"index"`             // 会計年度終了日
    
    // === 生データ（JSON保存） ===
    RawJSON             string `gorm:"type:jsonb;not null"`      // J-Quantsレスポンスそのまま
}
```

### 書類種別（TypeOfDocument）の例
| 値 | 説明 |
|----|------|
| FYFinancialStatements_Consolidated_JP | 決算短信 (連結・日本基準) |
| 1QFinancialStatements_Consolidated_JP | 第1四半期決算短信 (連結・日本基準) |
| 2QFinancialStatements_Consolidated_JP | 第2四半期決算短信 (連結・日本基準) |
| 3QFinancialStatements_Consolidated_JP | 第3四半期決算短信 (連結・日本基準) |
| FYFinancialStatements_Consolidated_IFRS | 決算短信 (連結・IFRS) |
| ... | その他多数 |

### 設計メモ
- `DisclosureNumber` は J-Quants が発行する一意の開示番号
- `RawJSON` にはAPIレスポンスの1レコード分をそのまま格納
- Python側では `json.loads(RawJSON)` → pandas で解析

---

## 4. news_searches (ニュース検索バッチ)

Tavily Search API で実行した検索バッチを管理する。
1つの銘柄に対して複数のクエリ（ビジネスモデル、決算短信等）を同時に実行し、
結果をまとめて保存する。

### Goモデル定義
```go
type NewsSearch struct {
    ID              uint      `gorm:"primaryKey"`
    Code            string    `gorm:"index;not null"`            // 銘柄コード
    SearchedAt      time.Time `gorm:"index;not null"`            // 検索実行日時
    CombinedContent string    `gorm:"type:text"`                 // まとめたテキスト（Pythonに渡す用）
    
    // リレーション
    Items           []NewsItem `gorm:"foreignKey:NewsSearchID"`
}
```

### 設計メモ
- `CombinedContent` には全クエリの結果を結合したテキストを保存
- gRPC で Python に渡す `qualitative_info` はこのフィールドから取得

---

## 5. news_items (ニュース検索結果)

Tavily Search API から返された個別の検索結果を保存する。

### APIレスポンス例
```go
type TavilySearchResponse struct {
    Results []TavilyResult `json:"results"`
    Answer  string         `json:"answer,omitempty"`
}

type TavilyResult struct {
    Title   string  `json:"title"`
    URL     string  `json:"url"`
    Content string  `json:"content"`
    Score   float64 `json:"score"`
}
```

### Goモデル定義
```go
type NewsItem struct {
    ID           uint    `gorm:"primaryKey"`
    NewsSearchID uint    `gorm:"index;not null"`              // 検索バッチへの外部キー
    SearchQuery  string  `gorm:"index"`                       // 検索クエリ（例: "楽天 ビジネスモデル"）
    Title        string
    URL          string
    Content      string  `gorm:"type:text"`
    Score        float64
    
    // リレーション
    NewsSearch   NewsSearch `gorm:"foreignKey:NewsSearchID"`
}
```

### 設計メモ
- `NewsSearchID` で同じタイミングで実行した検索をグルーピング
- `SearchQuery` で「どのクエリの結果か」を識別
- 同じURLが複数クエリで返ってくる可能性があるため、URLにユニーク制約は付けない

### 検索クエリの例
```
- "{企業名} ビジネスモデル"
- "{企業名} 決算短信 要約"
```

---

## 6. analysis_results (AI分析結果)

Python (LangGraph) からの分析結果を保存する。
gRPC Response の内容を永続化。

### Goモデル定義
```go
type AnalysisResult struct {
    ID               uint      `gorm:"primaryKey"`
    Code             string    `gorm:"index;not null"`           // 銘柄コード
    AnalyzedAt       time.Time `gorm:"index;not null"`           // 分析日時
    Sentiment        string    `gorm:"index"`                    // 投資判断（Strong Buy/Buy/Hold/Sell）
    SummaryReasoning string    `gorm:"type:text"`                // 分析レポート（Markdown）
    ThoughtLog       string    `gorm:"type:text"`                // AI思考ログ（デバッグ用）
    
    // 分析に使用したデータの参照（スナップショット情報）
    PriceDataFrom    string                                      // 株価データの開始日
    PriceDataTo      string                                      // 株価データの終了日
    NewsSearchID     *uint     `gorm:"index"`                    // 使用したニュース検索バッチ（nullable）
}
```

### Sentimentの値
| 値 | 説明 |
|----|------|
| Strong Buy | 強い買い推奨 |
| Buy | 買い推奨 |
| Hold | 保持推奨 |
| Sell | 売り推奨 |

### 設計メモ
- `PriceDataFrom` / `PriceDataTo` で「どの期間の株価データを使ったか」を記録
- `NewsSearchID` で「どのニュース検索結果を使ったか」を紐付け
- 同じ銘柄でも複数回分析できる（時系列での分析結果を追跡可能）

---

## ER図

```
┌─────────────────┐
│    companies    │
│─────────────────│
│ ID (PK)         │
│ Code (UK)       │◄──────────────────────────────────────┐
│ CompanyName     │                                       │
│ Sector17Code    │                                       │
│ ...             │                                       │
└─────────────────┘                                       │
                                                          │
┌─────────────────┐                                       │
│  daily_quotes   │                                       │
│─────────────────│                                       │
│ ID (PK)         │                                       │
│ Code ───────────┼───────────────────────────────────────┤
│ Date            │                                       │
│ Open/High/Low/  │                                       │
│ Close/Volume    │                                       │
│ ...             │                                       │
└─────────────────┘                                       │
                                                          │
┌─────────────────┐                                       │
│ financial_      │                                       │
│ statements      │                                       │
│─────────────────│                                       │
│ ID (PK)         │                                       │
│ Code ───────────┼───────────────────────────────────────┤
│ DisclosureNum   │                                       │
│ (UK)            │                                       │
│ RawJSON         │                                       │
│ ...             │                                       │
└─────────────────┘                                       │
                                                          │
┌─────────────────┐       ┌─────────────────┐            │
│  news_searches  │       │   news_items    │            │
│─────────────────│       │─────────────────│            │
│ ID (PK)         │◄──────│ NewsSearchID(FK)│            │
│ Code ───────────┼───────┼─────────────────┼────────────┤
│ SearchedAt      │       │ SearchQuery     │            │
│ CombinedContent │       │ Title/URL       │            │
└─────────────────┘       │ Content/Score   │            │
        ▲                 └─────────────────┘            │
        │                                                 │
        │                                                 │
┌───────┴─────────┐                                       │
│analysis_results │                                       │
│─────────────────│                                       │
│ ID (PK)         │                                       │
│ Code ───────────┼───────────────────────────────────────┘
│ AnalyzedAt      │
│ Sentiment       │
│ SummaryReasoning│
│ ThoughtLog      │
│ NewsSearchID(FK)│
└─────────────────┘

※ Code は外部キー制約ではなく、論理的な紐付け
  （J-Quants のデータ構造に合わせて文字列で管理）
```

---

## 次のステップ

1. ✅ DB設計完了
2. ⬜ `backend/models/japanese_stock.go` にGoモデルを実装
3. ⬜ マイグレーション実行
4. ⬜ Repository層の実装
