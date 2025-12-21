package models

import (
	"time"

	"gorm.io/gorm"
)

// Company 銘柄マスタ
// J-Quants /listed/info から取得した銘柄情報を保存する
type Company struct {
	gorm.Model
	Code               string `gorm:"uniqueIndex;not null" json:"Code"` // 銘柄コード（例: "86970"）
	CompanyName        string `json:"CompanyName"`                      // 企業名
	CompanyNameEnglish string `json:"CompanyNameEnglish"`               // 企業名（英語）
	Sector17Code       string `gorm:"index" json:"Sector17Code"`        // 17業種コード
	Sector17CodeName   string `json:"Sector17CodeName"`                 // 17業種名
	Sector33Code       string `gorm:"index" json:"Sector33Code"`        // 33業種コード
	Sector33CodeName   string `json:"Sector33CodeName"`                 // 33業種名
	ScaleCategory      string `json:"ScaleCategory"`                    // 規模区分
	MarketCode         string `gorm:"index" json:"MarketCode"`          // 市場コード
	MarketCodeName     string `json:"MarketCodeName"`                   // 市場名
}

// DailyQuote 日足株価
// J-Quants /prices/daily_quotes から取得した日足株価データを保存する（過去6ヶ月分）
type DailyQuote struct {
	ID               uint    `gorm:"primaryKey" json:"ID"`
	Code             string  `gorm:"index:idx_daily_quote_code_date,unique;not null" json:"Code"`
	Date             string  `gorm:"index:idx_daily_quote_code_date,unique;not null" json:"Date"`
	Open             float64 `json:"Open"`
	High             float64 `json:"High"`
	Low              float64 `json:"Low"`
	Close            float64 `json:"Close"`
	Volume           float64 `json:"Volume"`
	TurnoverValue    float64 `json:"TurnoverValue"`    // 売買代金
	AdjustmentFactor float64 `json:"AdjustmentFactor"` // 調整係数
	AdjustmentOpen   float64 `json:"AdjustmentOpen"`   // 調整後始値
	AdjustmentHigh   float64 `json:"AdjustmentHigh"`   // 調整後高値
	AdjustmentLow    float64 `json:"AdjustmentLow"`    // 調整後安値
	AdjustmentClose  float64 `json:"AdjustmentClose"`  // 調整後終値
	AdjustmentVolume float64 `json:"AdjustmentVolume"` // 調整後出来高
}

// FinancialStatement 財務諸表
// J-Quants /fins/statements から取得した財務諸表データを保存する（過去5年分）
// ハイブリッド方式: 検索・フィルタ用フィールドのみカラム化、残りはRawJSONとして保存
type FinancialStatement struct {
	ID                         uint   `gorm:"primaryKey" json:"ID"`
	Code                       string `gorm:"index;not null" json:"Code"`                   // 銘柄コード
	DisclosureNumber           string `gorm:"uniqueIndex;not null" json:"DisclosureNumber"` // 開示番号（一意識別子）
	DisclosedDate              string `gorm:"index" json:"DisclosedDate"`                   // 開示日
	TypeOfDocument             string `gorm:"index" json:"TypeOfDocument"`                  // 書類種別
	TypeOfCurrentPeriod        string `json:"TypeOfCurrentPeriod"`                          // 期種別（1Q/2Q/3Q/FY）
	CurrentFiscalYearStartDate string `json:"CurrentFiscalYearStartDate"`                   // 会計年度開始日
	CurrentFiscalYearEndDate   string `gorm:"index" json:"CurrentFiscalYearEndDate"`        // 会計年度終了日
	RawJSON                    string `gorm:"type:jsonb;not null" json:"RawJSON"`           // J-Quantsレスポンスそのまま
}

// NewsSearch ニュース検索バッチ
// Tavily Search API で実行した検索バッチを管理する
// 1つの銘柄に対して複数のクエリ（ビジネスモデル、決算短信等）を同時に実行し、結果をまとめて保存する
type NewsSearch struct {
	gorm.Model
	Code            string    `gorm:"index;not null" json:"Code"`       // 銘柄コード
	SearchedAt      time.Time `gorm:"index;not null" json:"SearchedAt"` // 検索実行日時
	CombinedContent string    `gorm:"type:text" json:"CombinedContent"` // まとめたテキスト（Pythonに渡す用）

	// リレーション
	Items []NewsItem `gorm:"foreignKey:NewsSearchID" json:"Items,omitempty"`
}

// NewsItem ニュース検索結果（個別）
// Tavily Search API から返された個別の検索結果を保存する
type NewsItem struct {
	gorm.Model
	NewsSearchID uint    `gorm:"index;not null" json:"NewsSearchID"` // 検索バッチへの外部キー
	SearchQuery  string  `gorm:"index" json:"SearchQuery"`           // 検索クエリ（例: "楽天 ビジネスモデル"）
	Title        string  `json:"Title"`
	URL          string  `json:"URL"`
	Content      string  `gorm:"type:text" json:"Content"`
	Score        float64 `json:"Score"`

	// リレーション
	NewsSearch NewsSearch `gorm:"foreignKey:NewsSearchID" json:"NewsSearch,omitempty"`
}

// AnalysisResult AI分析結果
// Python (LangGraph) からの分析結果を保存する
// gRPC Response の内容を永続化
type AnalysisResult struct {
	gorm.Model
	// Phase1, 2に共通する基本的な情報
	Code             string    `gorm:"index;not null" json:"Code"`        // 銘柄コード
	AnalyzedAt       time.Time `gorm:"index;not null" json:"AnalyzedAt"`  // 分析日時

	// Phase 2の推論結果
	StockSummary     string    `gorm:"type:text" json:"stocksummary"`     // AIによる企業の株価データの要約
	FinancialSummary string    `gorm:"type:text" json:"financialsummary"` // AIによる企業の財務情報の要約
	Sentiment        string    `gorm:"index" json:"Sentiment"`            // 投資判断（Strong Buy/Buy/Hold/Sell）
	SummaryReasoning string    `gorm:"type:text" json:"SummaryReasoning"` // 分析レポート（Markdown）
	ThoughtLog       string    `gorm:"type:text" json:"ThoughtLog"`       // AI思考ログ（デバッグ用）

	// Phase 1の推論結果
	BusinessModel string `gorm:"type:text" json:"BusinessModel"` // Phase1でAIが判断したその企業のビジネスモデル
	KPI           string `gorm:"type:text" json:"KPI"` // Phase1でAIが判断したその企業の見るべき項目

	// 分析に使用したデータの参照（スナップショット情報）
	PriceDataFrom string `json:"PriceDataFrom"`                       // 株価データの開始日
	PriceDataTo   string `json:"PriceDataTo"`                         // 株価データの終了日
	NewsSearchID  *uint  `gorm:"index" json:"NewsSearchID,omitempty"` // 使用したニュース検索バッチ（nullable）

	// リレーション
	SectorAnalysisResultID *uint `gorm:"index" json:"SectorAnalysisResultID,omitempty"`
	SectorAnalysisResult *SectorAnalysisResult `gorm:"foreignKey:SectorAnalysisResultID" json:"SectorAnalysisResult,omitempty"`
}

type SectorAnalysisResult struct {
	gorm.Model
	SectorCode      string `gorm:"index;not null" json:"SectorCode"` // セクターコード
	AnalyzedAt      time.Time `gorm:"index;not null" json:"AnalyzedAt"` // 分析日時
	
	// Top3銘柄の情報
	Top1Code        string `gorm:"index" json:"Top1Code"` // 1位の銘柄コード
	Top1Reasoning   string `gorm:"type:text" json:"Top1Reasoning"` // 1位の選出理由
	Top2Code        string `gorm:"index" json:"Top2Code"` // 2位の銘柄コード
	Top2Reasoning   string `gorm:"type:text" json:"Top2Reasoning"` // 2位の選出理由
	Top3Code        string `gorm:"index" json:"Top3Code"` // 3位の銘柄コード
	Top3Reasoning   string `gorm:"type:text" json:"Top3Reasoning"` // 3位の選出理由

	// 比較分析を行った際の思考ログ
	ComparisonLog   *string `gorm:"type:text" json:"ComparisonLog"`   // 比較分析の思考ログ 
	OverallSummary *string `gorm:"type:text" json:"OverallSummary"` //セクター分析のサマリ（投資の順位と共に出させてもいいけど、負荷が重そうなので、現状は出させない）

	// リレーション
	AnalysisResults []AnalysisResult `gorm:"foreignKey:SectorAnalysisResultID" json:"AnalysisResultID,omitempty"`
}



