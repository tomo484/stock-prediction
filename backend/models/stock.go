package models

import "gorm.io/gorm"

type DailyRanking struct {
	gorm.Model
	StockID      uint    `json:"StockID"`                      // Foreign Key (Stockテーブルへの紐付け)
	Date         string  `gorm:"index" json:"Date"`            // 日付 (例: "2025-01-01") ※米国基準（米国時間）
	Rank         int     `json:"Rank"`                         // その日の順位
	Category     string  `gorm:"index" json:"Category"`        // カテゴリ（例: Technology, Finance, Healthcare, etc.）
	ChangeAmount float64 `json:"ChangeAmount"`                 // 上昇額
	ChangeRate   float64 `json:"ChangeRate"`                   // 上昇率
	Price        float64 `json:"Price"`                        // その時の株価
	NewsSummary  string  `gorm:"type:text" json:"NewsSummary"` // AIに読ませたニュースの要約（念のため保存）
	AiAnalysis   string  `gorm:"type:text" json:"AiAnalysis"`  // AIが出した上昇理由
	Stock        Stock   `json:"Stock,omitempty"`
}

type Stock struct {
	gorm.Model
	Ticker   string         `gorm:"uniqueIndex;not null" json:"Ticker"` // 銘柄コード (例: NVDA)
	Name     string         `json:"Name"`// 企業名
	Sector   string         `json:"Sector"` //セクター（大分類）
	Industry string         `json:"Industry"` //業界（小分類）
	Description string      `gorm:"type:text"json:"Description"` //企業の説明
	Website string         `json:"Website"` //企業の公式WebsiteURL
	Country string         `json:"Country"` //企業の本社所在地
    FullTimeEmployees int  `json:"FullTimeEmployees"` //企業の従業員数
    Image             string `json:"Image"` //企業のロゴ画像URL
    IpoDate           string `json:"IpoDate"` //企業のIPO日
    CEO               string `json:"CEO"` //企業のCEO

	// リレーション
	Ranking  []DailyRanking `gorm:"foreignKey:StockID" json:"DailyRanking,omitempty"`
	Metrics  []StockMetric  `gorm:"foreignKey:StockID" json:"StockMetric,omitempty"`
}

type StockMetric struct {
    gorm.Model
    StockID       uint    `gorm:"index"` // Foreign Key (Stockテーブルへの紐付け)
    Date          string  `gorm:"index"` // 日付
    MarketCap     float64 `json:"MarketCap"`// 時価総額
    Volume        int64 `json:"Volume"` // 本日の出来高
    AverageVolume int64 `json:"AverageVolume"` //過去3か月平均出来高
    Beta          float64 `json:"Beta"` //ベータ値（市場全体に対する株価の感応度）
    LastDividend  float64 `json:"LastDividend"` //直近の一株当たりの配当金額
    
	// リレーション
    Stock         Stock   `json:"Stock,omitempty"`
}
