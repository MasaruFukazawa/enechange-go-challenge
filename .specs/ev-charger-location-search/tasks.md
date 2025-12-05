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
- **依存関係**: なし
- **ファイル**:
  - `models/location.go`
  - `models/evse.go`

### CORE-002: CSVデータインポート機能

- **説明**: サンプルCSVからデータベースへのインポート
- **受入条件**:
  - locations.csv がインポートできる
  - evses.csv がインポートできる
- **依存関係**: CORE-001
- **ファイル**:
  - `database/seed.go`

---

## APIタスク (API)

### API-001: ロケーション検索エンドポイント実装

- **説明**: GET /api/locations エンドポイントの実装
- **受入条件**:
  - latitude, longitude パラメータが必須
  - radius パラメータがオプション（デフォルト100km）
  - date_from, date_to パラメータがオプション
  - 円形範囲内のロケーションが返却される
- **依存関係**: CORE-001, CORE-002
- **ファイル**:
  - `router/router.go`
  - `handlers/location_handler.go`
  - `services/location_service.go`

---

## タスク依存関係図

```
【フェーズ1: データ層】
CORE-001 (モデル定義)
    ↓
CORE-002 (CSVインポート)

【フェーズ2: API層】
CORE-002 → API-001 (エンドポイント実装)
```
