package xpost

import (
	"fmt"
	"stock-prediction/backend/models"
	"strings"
	"time"
)

// ランキング（AIAnalysis無し）投稿用テンプレート
func BuildRankingPost(date string, rankings []models.DailyRanking) string {
	//日付フォーマット："2025-11-27" -> "11/27"
	t, _ := time.Parse("2006-01-02", date)
	dateStr := t.Format("1/2")

	models := []string{"🥇", "🥈", "🥉", "4️⃣", "5️⃣"}

	var lines []string
	lines = append(lines, fmt.Sprintf("🚀 %s 米国株急騰ランキング", dateStr))

	for i, ranking := range rankings {
		if i >= 5 {
			break
		}
		lines = append(lines, fmt.Sprintf("%s %s (+%.1f%%)", models[i], ranking.Stock.Name, ranking.ChangeRate))
	}

	return strings.Join(lines, "\n")
}

// 個別分析（AIAnalysisあり）投稿用テンプレート
func BuildAnalysisPost(ranking models.DailyRanking) string {
	medals := []string{"🥇", "🥈", "🥉", "4️⃣", "5️⃣"}
	medal := medals[ranking.Rank-1]

	header := fmt.Sprintf("%s %s (+%.1f%%)\n",
		medal, ranking.Stock.Name, ranking.ChangeRate)

	maxAnalysisLen := 280 - len([]rune(header))
	truncatedAnalysis := TruncateForX(ranking.AiAnalysis, maxAnalysisLen)

	return header + truncatedAnalysis
}

// AiAnalysisを安全な文字数に切り詰める
func TruncateForX(text string, maxLen int) string {
	if len([]rune(text)) <= maxLen {
		return text
	}
	// 末尾に「...」を追加
	runes := []rune(text)
	return string(runes[:maxLen-3]) + "..."
}
