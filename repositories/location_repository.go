package repositories

import (
	"go-challenge/models"

	"gorm.io/gorm"
)

// LocationRepository はLocationエンティティのデータアクセスインターフェース
type LocationRepository interface {
	FindAll() ([]models.Location, error)
	FindAllWithEVSEs() ([]models.Location, error)
}

type locationRepository struct {
	db *gorm.DB
}

// NewLocationRepository は新しいLocationRepositoryを作成する
func NewLocationRepository(db *gorm.DB) LocationRepository {
	return &locationRepository{db: db}
}

// FindAll は全ロケーションを取得する
func (r *locationRepository) FindAll() ([]models.Location, error) {
	var locations []models.Location
	if err := r.db.Find(&locations).Error; err != nil {
		return nil, err
	}
	return locations, nil
}

// FindAllWithEVSEs はEVSE含む全ロケーションを取得する
func (r *locationRepository) FindAllWithEVSEs() ([]models.Location, error) {
	var locations []models.Location
	if err := r.db.Preload("EVSEs").Find(&locations).Error; err != nil {
		return nil, err
	}
	return locations, nil
}
