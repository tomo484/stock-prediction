package japanese_Stock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"stock-prediction/backend/models"
	"stock-prediction/backend/repositories"
	"stock-prediction/backend/utils"
	"time"
)

// FinancialStatementResponse J-Quants APIレスポンス用の型（DBモデルとは別）
// APIレスポンス全体をmapとして受け取り、必要なフィールドを抽出してRawJSONとして保存
type FinancialStatementResponse map[string]interface{}

type ListedFinancialStatementsResponse struct {
	FinancialInfo []FinancialStatementResponse `json:"financial_info"`
}

// FetchJQuantsFinancialStatements 財務諸表データを取得する
func FetchJQuantsFinancialStatements(idToken string, code string) (*ListedFinancialStatementsResponse, error) {
	url := fmt.Sprintf("https://api.jquants.com/v1/fins/statements?code=%s", code)
	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// IDトークンをヘッダーに付与
	req.Header.Set("Authorization", "Bearer "+idToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch financial statements: status %d", resp.StatusCode)
	}

	var result ListedFinancialStatementsResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// SaveJQuantsFinancialStatementsToDB 財務諸表データをDBに保存する
func SaveJQuantsFinancialStatementsToDB(financialStatements *ListedFinancialStatementsResponse, repository repositories.IJapaneseStockRepository) error {
	for _, fs := range financialStatements.FinancialInfo {
		// APIレスポンス全体をJSON文字列に変換してRawJSONとして保存
		rawJSONBytes, err := json.Marshal(fs)
		if err != nil {
			return fmt.Errorf("failed to marshal raw JSON: %w", err)
		}

		// mapから必要なフィールドを抽出（utils関数を使用して安全に取得）
		code := utils.GetStringFromMap(fs, "Code")
		disclosureNumber := utils.GetStringFromMap(fs, "DisclosureNumber")
		disclosedDate := utils.GetStringFromMap(fs, "DisclosedDate")
		typeOfDocument := utils.GetStringFromMap(fs, "TypeOfDocument")
		typeOfCurrentPeriod := utils.GetStringFromMap(fs, "TypeOfCurrentPeriod")
		currentFiscalYearStartDate := utils.GetStringFromMap(fs, "CurrentFiscalYearStartDate")
		currentFiscalYearEndDate := utils.GetStringFromMap(fs, "CurrentFiscalYearEndDate")

		// FinancialStatementResponse -> models.FinancialStatementに変換する
		financialStatement := &models.FinancialStatement{
			Code:                       code,
			DisclosureNumber:           disclosureNumber,
			DisclosedDate:              disclosedDate,
			TypeOfDocument:             typeOfDocument,
			TypeOfCurrentPeriod:        typeOfCurrentPeriod,
			CurrentFiscalYearStartDate: currentFiscalYearStartDate,
			CurrentFiscalYearEndDate:   currentFiscalYearEndDate,
			RawJSON:                    string(rawJSONBytes),
		}

		// 既存なら更新、新規なら保存する
		err = repository.CreateOrUpdateFinancialStatement(financialStatement)
		if err != nil {
			return fmt.Errorf("failed to create/update financial statement %s %s: %w", code, disclosureNumber, err)
		}
	}

	return nil
}

func SyncJQuantsFinancialStatements(idToken string, code string, repository repositories.IJapaneseStockRepository) ([]models.FinancialStatement, error) {
	// JQuants APIから財務諸表データを取得
	financialStatements, err := FetchJQuantsFinancialStatements(idToken, code)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JQuants financial statements: %w", err)
	}

	// 財務諸表データをDBに保存
	err = SaveJQuantsFinancialStatementsToDB(financialStatements, repository)
	if err != nil {
		return nil, fmt.Errorf("failed to save JQuants financial statements to DB: %w", err)
	}

	// 取得し、保存した財務諸表データを返す（型の整合性とデータの正確性を保証）
	savedStatements, err := repository.FindFinancialStatementsByCode(code)
	if err != nil {
		return nil, fmt.Errorf("failed to find financial statements by code: %w", err)
	}

	return savedStatements, nil
}
