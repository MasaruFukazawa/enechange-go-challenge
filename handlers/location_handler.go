package handlers

import (
	"net/http"
	"regexp"
	"strconv"
	"time"

	"go-challenge/services"

	"github.com/gin-gonic/gin"
)

var (
	latitudeRegex  = regexp.MustCompile(`^-?[0-9]{1,2}\.[0-9]{5,7}$`)
	longitudeRegex = regexp.MustCompile(`^-?[0-9]{1,3}\.[0-9]{5,7}$`)
)

// ErrorResponse はエラーレスポンス
type ErrorResponse struct {
	Error string `json:"error"`
}

// LocationHandler はHTTPリクエストを処理するハンドラー
type LocationHandler struct {
	service services.LocationService
}

// NewLocationHandler は新しいLocationHandlerを作成する
func NewLocationHandler(service services.LocationService) *LocationHandler {
	return &LocationHandler{service: service}
}

// GetLocations はロケーション検索エンドポイント
func (h *LocationHandler) GetLocations(c *gin.Context) {
	// パラメータ取得
	latStr := c.Query("latitude")
	lonStr := c.Query("longitude")
	radiusStr := c.Query("radius")
	dateFromStr := c.Query("date_from")
	dateToStr := c.Query("date_to")

	// 必須パラメータチェック
	if latStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "latitude is required"})
		return
	}
	if lonStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "longitude is required"})
		return
	}

	// フォーマットバリデーション
	if !latitudeRegex.MatchString(latStr) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid latitude format"})
		return
	}
	if !longitudeRegex.MatchString(lonStr) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid longitude format"})
		return
	}

	// パース
	lat, _ := strconv.ParseFloat(latStr, 64)
	lon, _ := strconv.ParseFloat(lonStr, 64)

	// 半径（デフォルト100km）
	radius := 100.0
	if radiusStr != "" {
		r, err := strconv.Atoi(radiusStr)
		if err != nil || r <= 0 {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid radius value"})
			return
		}
		radius = float64(r)
	}

	// 日時パラメータ
	var dateFrom, dateTo *time.Time
	if dateFromStr != "" {
		t, err := time.Parse(time.RFC3339, dateFromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid date_from format"})
			return
		}
		dateFrom = &t
	}
	if dateToStr != "" {
		t, err := time.Parse(time.RFC3339, dateToStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid date_to format"})
			return
		}
		dateTo = &t
	}

	// 検索実行
	params := services.SearchParams{
		Latitude:  lat,
		Longitude: lon,
		Radius:    radius,
		DateFrom:  dateFrom,
		DateTo:    dateTo,
	}

	results, err := h.service.SearchByRadius(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, results)
}
