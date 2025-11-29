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
	Name     string         `json:"Name"`
	Sector   string         `json:"Sector"`
	Industry string         `json:"Industry"`
	Ranking  []DailyRanking `gorm:"foreignKey:StockID" json:"Ranking,omitempty"`
}
