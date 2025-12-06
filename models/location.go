package models

import (
	"fmt"
	"time"
)

// Location はロケーションエンティティ
type Location struct {
	ID          string     `gorm:"primaryKey;size:36" json:"id"`
	Name        *string    `gorm:"size:255" json:"name"`
	Address     string     `gorm:"size:45;not null" json:"address"`
	Latitude    float64    `gorm:"not null" json:"-"`
	Longitude   float64    `gorm:"not null" json:"-"`
	LastUpdated *time.Time `json:"-"`
	EVSEs       []EVSE     `gorm:"foreignKey:LocationID" json:"evses"`
}

// TableName はGORMのテーブル名を返す
func (Location) TableName() string {
	return "locations"
}

// LocationResponse はJSON出力用の構造体
type LocationResponse struct {
	ID          string         `json:"id"`
	Name        *string        `json:"name,omitempty"`
	Address     string         `json:"address"`
	Coordinates CoordinatesDTO `json:"coordinates"`
	EVSEs       []EVSEResponse `json:"evses"`
}

// CoordinatesDTO は座標情報のDTO
type CoordinatesDTO struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

// ToResponse はLocationをJSON出力用構造体に変換する
func (l Location) ToResponse() LocationResponse {
	evseResponses := make([]EVSEResponse, len(l.EVSEs))
	for i, evse := range l.EVSEs {
		evseResponses[i] = evse.ToResponse()
	}

	return LocationResponse{
		ID:      l.ID,
		Name:    l.Name,
		Address: l.Address,
		Coordinates: CoordinatesDTO{
			Latitude:  fmt.Sprintf("%g", l.Latitude),
			Longitude: fmt.Sprintf("%g", l.Longitude),
		},
		EVSEs: evseResponses,
	}
}
