package services

import (
	"os"
	"github.com/sashabaranov/go-openai"
	"fmt"
	"context"
)

func AnalyzeStockRise(ticker string, changeRate float64, newsHeadlines []string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)

	systemPrompt := `
あなたはプロの株式市場アナリストです。
提供された「銘柄」「上昇率」「関連ニュース」をもとに、
その株がなぜ急上昇したのか、その要因を簡潔に日本語で解説してください。
ニュースがない場合は、その企業の一般的な事業内容と、この上昇率が通常あり得るものかどうかを述べてください。
回答は150文字以内で、投資家向けに要約してください。
`
    newsText := "特になし"
	if len(newsHeadlines) > 0 {
		newsText = ""
		for _, h := range newsHeadlines {
			newsText += "- " + h + "\n"
		}
	}

	userContent := fmt.Sprintf(
		"銘柄: %s\n本日の上昇率: +%.2f%%\n関連ニュース:\n%s\n\nこの上昇の理由を分析してください。",
		ticker, changeRate, newsText,
	)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "gpt-4o", // または "gpt-3.5-turbo" (安い)
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userContent,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

