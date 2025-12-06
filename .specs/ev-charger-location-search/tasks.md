# 実装タスク: ev-charger-location-search

## タスク概要

EV充電器ロケーション検索APIの実装

---

## コアタスク (CORE)

### CORE-001: データベースモデル定義

- **説明**: Location と EVSE のGORMモデルを定義
- **受入条件**:
  - Location モデルが定義されている
  - EVSE モデルが定義されている
  - リレーションが正しく設定されている
  - ToResponse メソッドのテストが通る
- **依存関係**: なし
- **ファイル**:
  - `models/location.go`
  - `models/evse.go`
  - `models/location_test.go`
  - `models/evse_test.go`

### CORE-002: CSVデータインポート機能

- **説明**: サンプルCSVからデータベースへのインポート
- **受入条件**:
  - locations.csv がインポートできる
  - evses.csv がインポートできる
  - インポート処理のテストが通る
- **依存関係**: CORE-001
- **ファイル**:
  - `database/seed.go`
  - `database/seed_test.go`

---

## APIタスク (API)

### API-001: ロケーション検索エンドポイント実装

- **説明**: GET /api/locations エンドポイントの実装
- **受入条件**:
  - latitude, longitude パラメータが必須
  - radius パラメータがオプション（デフォルト100km）
  - date_from, date_to パラメータがオプション
  - 円形範囲内のロケーションが返却される
  - 各レイヤーのユニットテストが通る
- **依存関係**: CORE-001, CORE-002
- **ファイル**:
  - `router/router.go`
  - `handlers/location_handler.go`
  - `handlers/location_handler_test.go`
  - `services/location_service.go`
  - `services/location_service_test.go`
  - `services/distance.go`
  - `services/distance_test.go`
  - `repositories/location_repository.go`
  - `repositories/location_repository_test.go`

---

## テストタスク (TEST)

### TEST-001: 統合テスト

- **説明**: E2Eでの検索機能テスト
- **受入条件**:
  - CSVデータからの検索が正しく動作する
  - レスポンスJSONの形式が正しい
  - EVSEが正しく含まれる
- **依存関係**: API-001
- **ファイル**:
  - `handlers/location_handler_integration_test.go`

---

## タスク依存関係図

```
【フェーズ1: データ層】
CORE-001 (モデル定義 + ユニットテスト)
    ↓
CORE-002 (CSVインポート + ユニットテスト)

【フェーズ2: API層】
CORE-002 → API-001 (エンドポイント実装 + ユニットテスト)

【フェーズ3: 統合テスト】
API-001 → TEST-001 (統合テスト)
```
