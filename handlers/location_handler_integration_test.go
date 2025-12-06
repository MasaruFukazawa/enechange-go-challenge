// +build integration

package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"go-challenge/database"
	"go-challenge/models"
	"go-challenge/repositories"
	"go-challenge/services"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupIntegrationTestDB(t *testing.T) *gorm.DB {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(go-challenge-mysql:3306)/go-challenge_development?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// テストデータをシード
	locationsPath := filepath.Join("..", "sample", "locations.csv")
	evsesPath := filepath.Join("..", "sample", "evses.csv")
	if err := database.SeedFromCSV(db, locationsPath, evsesPath); err != nil {
		t.Fatalf("Failed to seed test data: %v", err)
	}

	return db
}

func setupIntegrationRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	locationRepo := repositories.NewLocationRepository(db)
	locationService := services.NewLocationService(locationRepo)
	locationHandler := NewLocationHandler(locationService)

	r.GET("/api/locations", locationHandler.GetLocations)
	return r
}

func TestIntegration_GetLocations_ReturnsLocationsWithEVSEs(t *testing.T) {
	db := setupIntegrationTestDB(t)
	router := setupIntegrationRouter(db)

	// 旭川付近で検索（半径50km）
	req, _ := http.NewRequest("GET", "/api/locations?latitude=43.796611&longitude=142.375917&radius=50", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []models.LocationResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) == 0 {
		t.Error("Expected non-empty response for Asahikawa area")
	}

	// EVSEが含まれていることを確認
	hasEVSE := false
	for _, loc := range response {
		if len(loc.EVSEs) > 0 {
			hasEVSE = true
			break
		}
	}
	if !hasEVSE {
		t.Error("Expected at least one location with EVSEs")
	}

	t.Logf("Found %d locations in 50km radius of Asahikawa", len(response))
}

func TestIntegration_GetLocations_EmptyResultForRemoteArea(t *testing.T) {
	db := setupIntegrationTestDB(t)
	router := setupIntegrationRouter(db)

	// 遠く離れた場所で検索（半径1km）
	req, _ := http.NewRequest("GET", "/api/locations?latitude=0.000000&longitude=0.000000&radius=1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []models.LocationResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) != 0 {
		t.Errorf("Expected empty response for remote area, got %d locations", len(response))
	}
}

func TestIntegration_GetLocations_JSONFormat(t *testing.T) {
	db := setupIntegrationTestDB(t)
	router := setupIntegrationRouter(db)

	req, _ := http.NewRequest("GET", "/api/locations?latitude=43.796611&longitude=142.375917&radius=100", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response []models.LocationResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) == 0 {
		t.Skip("No locations found for format test")
	}

	loc := response[0]

	// 必須フィールドの確認
	if loc.ID == "" {
		t.Error("Location ID is empty")
	}
	if loc.Address == "" {
		t.Error("Location Address is empty")
	}
	if loc.Coordinates.Latitude == "" {
		t.Error("Location Coordinates.Latitude is empty")
	}
	if loc.Coordinates.Longitude == "" {
		t.Error("Location Coordinates.Longitude is empty")
	}

	// EVSEのステータスフォーマット確認
	if len(loc.EVSEs) > 0 {
		evse := loc.EVSEs[0]
		validStatuses := map[string]bool{
			"AVAILABLE": true, "BLOCKED": true, "CHARGING": true,
			"INOPERATIVE": true, "OUTOFORDER": true, "PLANNED": true,
			"REMOVED": true, "RESERVED": true, "UNKNOWN": true,
		}
		if !validStatuses[evse.Status] {
			t.Errorf("Invalid EVSE status: %s", evse.Status)
		}
	}
}
