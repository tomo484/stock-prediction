package japanese_Stock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"stock-prediction/backend/models"
	"stock-prediction/backend/repositories"
)

// CompanyInfo J-Quants /listed/info のレスポンス構造
type CompanyInfo struct {
	Date                string `json:"Date"`
	Code                string `json:"Code"`
	CompanyName         string `json:"CompanyName"`
	CompanyNameEnglish  string `json:"CompanyNameEnglish"`
	Sector17Code        string `json:"Sector17Code"`
	Sector17CodeName    string `json:"Sector17CodeName"`
	Sector33Code        string `json:"Sector33Code"`
	Sector33CodeName    string `json:"Sector33CodeName"`
	ScaleCategory       string `json:"ScaleCategory"`
	MarketCode          string `json:"MarketCode"`
	MarketCodeName      string `json:"MarketCodeName"`
}

type ListedInfoResponse struct {
	Info []CompanyInfo `json:"info"`
}

// 東証市場の銘柄マスタを取得する
func FetchJQuantsCompanies(idToken string) (*ListedInfoResponse, error) {
	url := "https://api.jquants.com/v1/listed/info"
	client := &http.Client{Timeout: 30 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// IDトークンをヘッダーに付与
	req.Header.Set("Authorization", "Bearer " +idToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch companies: status %d", resp.StatusCode)
	}

	var result ListedInfoResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// 東証市場の銘柄マスタをDBに保存する
func SaveJQuantsCompaniesToDB(companies *ListedInfoResponse, repository repositories.IJapaneseStockRepository) error {
	for _, company := range companies.Info {
		// CompanyInfo -> models.Companyに変換する
		company := &models.Company{
			Code:               company.Code,
			CompanyName:        company.CompanyName,
			CompanyNameEnglish: company.CompanyNameEnglish,
			Sector17Code:       company.Sector17Code,
			Sector17CodeName:   company.Sector17CodeName,
			Sector33Code:       company.Sector33Code,
			Sector33CodeName:   company.Sector33CodeName,
			ScaleCategory:      company.ScaleCategory,
			MarketCode:         company.MarketCode,
			MarketCodeName:     company.MarketCodeName,
		}

		// 既存なら更新、新規なら保存する
		err := repository.CreateOrUpdateCompany(company)
		if err != nil {
			return fmt.Errorf("failed to create/update company %s: %w", company.Code, err)
		}
	}

	return nil
}