package controllers

import (
	"fmt"
	"net/http"
	"stock-prediction/backend/repositories"
	"stock-prediction/backend/services"
	xpost "stock-prediction/backend/services/x_post"

	"github.com/labstack/echo/v4"
)

type IStockController interface {
	FindLatestRanking(c echo.Context) error
	FindDailyRanking(c echo.Context) error
	FindStock(c echo.Context) error
	SyncData(c echo.Context) error
	XAutomaticallyPost(c echo.Context) error
}

type stockController struct {
	service      services.IStockService
	xPostService xpost.IXPostService
}

func NewStockController(service services.IStockService, repo repositories.IStockRepository) IStockController {
	xPostService := xpost.NewXPostService(repo)
	return &stockController{service: service, xPostService: xPostService}
}

func (sc *stockController) FindLatestRanking(c echo.Context) error {
	latestRanking, err := sc.service.FindLatestRanking()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, latestRanking)
}

func (sc *stockController) FindDailyRanking(c echo.Context) error {
	date := c.QueryParam("date")
	if date == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "date query parameter is required"})
	}
	dailyRanking, err := sc.service.FindDailyRanking(date)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dailyRanking)
}

func (sc *stockController) FindStock(c echo.Context) error {
	ticker := c.Param("ticker")
	stock, err := sc.service.FindStock(ticker)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, stock)
}

func (sc *stockController) SyncData(c echo.Context) error {
	err := sc.service.SyncData()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Data synchronized successfully")
}

func (sc *stockController) XAutomaticallyPost(c echo.Context) error {
	// ===== 一時的なデバッグコード（削除予定） =====
	fmt.Printf("🔍 [DEBUG] XAutomaticallyPost 開始\n")
	// ===== デバッグコード終了 =====

	posttype := c.QueryParam("posttype")
	date := c.QueryParam("date")

	fmt.Printf("🔍 [DEBUG] パラメータ: posttype=%s, date=%s\n", posttype, date)

	if posttype == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "type query parameter is required",
		})
	}

	var err error
	var message string

	switch posttype {
	case "ranking":
		//ランキング投稿（AiAnalysis無し）
		fmt.Printf("🔍 [DEBUG] ランキング投稿を開始\n")
		err = sc.xPostService.PostRanking(date)
		message = "Ranking posted to X successfully"
	case "analysis":
		//個別分析投稿（5件まとめて）
		fmt.Printf("🔍 [DEBUG] 分析投稿を開始\n")
		err = sc.xPostService.PostAnalysis(date)
		message = "Analysis posted to X successfully"
	case "all":
		//ランキングと個別分析をまとめて投稿
		fmt.Printf("🔍 [DEBUG] ランキング+分析投稿を開始\n")
		err = sc.xPostService.PostRanking(date)
		if err != nil {
			fmt.Printf("🔍 [DEBUG] ランキング投稿エラー: %v\n", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "failed to post ranking to x:" + err.Error(),
			})
		}
		err = sc.xPostService.PostAnalysis(date)
		if err != nil {
			fmt.Printf("🔍 [DEBUG] 分析投稿エラー: %v\n", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "failed to post analysis to x:" + err.Error(),
			})
		}
		message = "Ranking and analysis posted to X successfully"
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid type. use 'ranking', 'analysis', or 'all'",
		})
	}
	if err != nil {
		fmt.Printf("🔍 [DEBUG] 投稿エラー: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to post to x:" + err.Error(),
		})
	}

	fmt.Printf("✅ [DEBUG] XAutomaticallyPost 成功\n")
	return c.JSON(http.StatusOK, map[string]string{
		"message": message,
	})
}
