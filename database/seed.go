package database

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"go-challenge/models"

	"gorm.io/gorm"
)

// SeedFromCSV はCSVファイルからデータベースに初期データを投入する
func SeedFromCSV(db *gorm.DB, locationsPath, evsesPath string) error {
	// テーブルの作成（AutoMigrate）
	if err := db.AutoMigrate(&models.Location{}, &models.EVSE{}); err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}

	// 既存データの削除（冪等性確保）
	if err := db.Exec("DELETE FROM evses").Error; err != nil {
		return fmt.Errorf("failed to delete evses: %w", err)
	}
	if err := db.Exec("DELETE FROM locations").Error; err != nil {
		return fmt.Errorf("failed to delete locations: %w", err)
	}

	// Locationsのインポート
	if err := importLocations(db, locationsPath); err != nil {
		return fmt.Errorf("failed to import locations: %w", err)
	}

	// EVSEsのインポート
	if err := importEVSEs(db, evsesPath); err != nil {
		return fmt.Errorf("failed to import evses: %w", err)
	}

	return nil
}

func importLocations(db *gorm.DB, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// ヘッダーをスキップ
	for i, record := range records {
		if i == 0 {
			continue // skip header
		}

		lat, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return fmt.Errorf("invalid latitude at row %d: %w", i, err)
		}

		lon, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return fmt.Errorf("invalid longitude at row %d: %w", i, err)
		}

		name := record[1]
		var namePtr *string
		if name != "" {
			namePtr = &name
		}

		location := models.Location{
			ID:        record[0],
			Name:      namePtr,
			Address:   record[2],
			Latitude:  lat,
			Longitude: lon,
		}

		if err := db.Create(&location).Error; err != nil {
			return fmt.Errorf("failed to create location at row %d: %w", i, err)
		}
	}

	return nil
}

func importEVSEs(db *gorm.DB, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// ヘッダーをスキップ
	for i, record := range records {
		if i == 0 {
			continue // skip header
		}

		status, err := strconv.Atoi(record[2])
		if err != nil {
			return fmt.Errorf("invalid status at row %d: %w", i, err)
		}

		evse := models.EVSE{
			LocationID: record[0],
			UID:        record[1],
			Status:     status,
		}

		if err := db.Create(&evse).Error; err != nil {
			return fmt.Errorf("failed to create evse at row %d: %w", i, err)
		}
	}

	return nil
}
