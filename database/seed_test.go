package database

import (
	"os"
	"path/filepath"
	"testing"

	"go-challenge/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(go-challenge-mysql:3306)/go-challenge_development?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// テーブルをクリーンアップ
	db.Exec("DROP TABLE IF EXISTS evses")
	db.Exec("DROP TABLE IF EXISTS locations")

	return db
}

func TestSeedFromCSV(t *testing.T) {
	db := setupTestDB(t)

	// サンプルデータのパスを取得
	locationsPath := filepath.Join("..", "sample", "locations.csv")
	evsesPath := filepath.Join("..", "sample", "evses.csv")

	// シード実行
	err := SeedFromCSV(db, locationsPath, evsesPath)
	if err != nil {
		t.Fatalf("SeedFromCSV failed: %v", err)
	}

	// Locationの件数を確認
	var locationCount int64
	db.Model(&models.Location{}).Count(&locationCount)
	if locationCount == 0 {
		t.Error("No locations were imported")
	}

	// EVSEの件数を確認
	var evseCount int64
	db.Model(&models.EVSE{}).Count(&evseCount)
	if evseCount == 0 {
		t.Error("No EVSEs were imported")
	}

	// 最初のロケーションの内容を確認
	var location models.Location
	db.First(&location)
	if location.ID == "" {
		t.Error("Location ID is empty")
	}
	if location.Address == "" {
		t.Error("Location Address is empty")
	}

	t.Logf("Imported %d locations and %d EVSEs", locationCount, evseCount)
}

func TestSeedFromCSV_Idempotent(t *testing.T) {
	db := setupTestDB(t)

	locationsPath := filepath.Join("..", "sample", "locations.csv")
	evsesPath := filepath.Join("..", "sample", "evses.csv")

	// 2回実行しても同じ結果になることを確認（冪等性）
	err := SeedFromCSV(db, locationsPath, evsesPath)
	if err != nil {
		t.Fatalf("First SeedFromCSV failed: %v", err)
	}

	var firstCount int64
	db.Model(&models.Location{}).Count(&firstCount)

	err = SeedFromCSV(db, locationsPath, evsesPath)
	if err != nil {
		t.Fatalf("Second SeedFromCSV failed: %v", err)
	}

	var secondCount int64
	db.Model(&models.Location{}).Count(&secondCount)

	if firstCount != secondCount {
		t.Errorf("Idempotency failed: first=%d, second=%d", firstCount, secondCount)
	}
}
