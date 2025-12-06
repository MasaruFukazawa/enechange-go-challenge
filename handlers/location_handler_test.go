package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-challenge/models"
	"go-challenge/services"

	"github.com/gin-gonic/gin"
)

// MockLocationService はテスト用のモックサービス
type MockLocationService struct{}

func (m *MockLocationService) SearchByRadius(params services.SearchParams) ([]models.LocationResponse, error) {
	name := "テストロケーション"
	return []models.LocationResponse{
		{
			ID:      "1",
			Name:    &name,
			Address: "東京都",
			Coordinates: models.CoordinatesDTO{
				Latitude:  "35.6762",
				Longitude: "139.6503",
			},
			EVSEs: []models.EVSEResponse{
				{UID: "evse-001", Status: "AVAILABLE"},
			},
		},
	}, nil
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := &MockLocationService{}
	handler := NewLocationHandler(mockService)
	r.GET("/api/locations", handler.GetLocations)
	return r
}

func TestGetLocations_Success(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/locations?latitude=35.676200&longitude=139.650300", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []models.LocationResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if len(response) == 0 {
		t.Error("Expected non-empty response")
	}
}

func TestGetLocations_MissingLatitude(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/locations?longitude=139.650300", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.Error != "latitude is required" {
		t.Errorf("Expected error 'latitude is required', got '%s'", response.Error)
	}
}

func TestGetLocations_MissingLongitude(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/locations?latitude=35.676200", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.Error != "longitude is required" {
		t.Errorf("Expected error 'longitude is required', got '%s'", response.Error)
	}
}

func TestGetLocations_InvalidLatitudeFormat(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/locations?latitude=invalid&longitude=139.650300", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.Error != "Invalid latitude format" {
		t.Errorf("Expected error 'Invalid latitude format', got '%s'", response.Error)
	}
}

func TestGetLocations_InvalidLongitudeFormat(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/locations?latitude=35.676200&longitude=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.Error != "Invalid longitude format" {
		t.Errorf("Expected error 'Invalid longitude format', got '%s'", response.Error)
	}
}

func TestGetLocations_InvalidRadius(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/locations?latitude=35.676200&longitude=139.650300&radius=-1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.Error != "Invalid radius value" {
		t.Errorf("Expected error 'Invalid radius value', got '%s'", response.Error)
	}
}

func TestGetLocations_WithValidRadius(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/locations?latitude=35.676200&longitude=139.650300&radius=50", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}
