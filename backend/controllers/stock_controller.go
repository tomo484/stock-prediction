package controllers

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"stock-prediction/backend/services"
)

type IStockController interface {
	FindLatestRanking(c echo.Context) error
	FindDailyRanking(c echo.Context) error
	FindStock(c echo.Context) error
	SyncData(c echo.Context) error
}

type stockController struct {
	service services.IStockService
}

func NewStockController(service services.IStockService) IStockController {
	return &stockController{service: service}
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