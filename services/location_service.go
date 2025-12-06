package services

import (
	"time"

	"go-challenge/models"
	"go-challenge/repositories"
)

// SearchParams は検索パラメータ
type SearchParams struct {
	Latitude  float64
	Longitude float64
	Radius    float64    // km
	DateFrom  *time.Time // inclusive
	DateTo    *time.Time // exclusive
}

// LocationService はロケーション検索のビジネスロジックインターフェース
type LocationService interface {
	SearchByRadius(params SearchParams) ([]models.LocationResponse, error)
}

type locationService struct {
	repo repositories.LocationRepository
}

// NewLocationService は新しいLocationServiceを作成する
func NewLocationService(repo repositories.LocationRepository) LocationService {
	return &locationService{repo: repo}
}

// SearchByRadius は指定された半径内のロケーションを検索する
func (s *locationService) SearchByRadius(params SearchParams) ([]models.LocationResponse, error) {
	// 全ロケーション（EVSE含む）を取得
	locations, err := s.repo.FindAllWithEVSEs()
	if err != nil {
		return nil, err
	}

	var results []models.LocationResponse

	for _, loc := range locations {
		// 距離を計算
		distance := HaversineDistance(params.Latitude, params.Longitude, loc.Latitude, loc.Longitude)

		// 半径外はスキップ
		if distance > params.Radius {
			continue
		}

		// 日時フィルター（DateFrom: inclusive）
		if params.DateFrom != nil && loc.LastUpdated != nil {
			if loc.LastUpdated.Before(*params.DateFrom) {
				continue
			}
		}

		// 日時フィルター（DateTo: exclusive）
		if params.DateTo != nil && loc.LastUpdated != nil {
			if !loc.LastUpdated.Before(*params.DateTo) {
				continue
			}
		}

		// レスポンス形式に変換
		results = append(results, loc.ToResponse())
	}

	// nilではなく空スライスを返す
	if results == nil {
		results = []models.LocationResponse{}
	}

	return results, nil
}
