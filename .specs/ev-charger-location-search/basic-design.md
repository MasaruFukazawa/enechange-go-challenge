# 基本設計: ev-charger-location-search

## 1. プロジェクト概要

### 1.1 目的と背景

ENECHANGE Go Developer Challengeの一環として、EV（電気自動車）充電器の検索APIを実装する。このAPIは、指定された緯度・経度・半径で定義された円形エリア内にある充電ロケーション（Location）とその充電設備（EVSE: Electric Vehicle Supply Equipment）を検索・返却する機能を提供する。

### 1.2 スコープ

- **含まれるもの**:

  - `GET /api/locations` エンドポイントの実装
  - 緯度・経度・半径による円形範囲検索機能
  - 更新日時によるフィルタリング機能（date_from, date_to）
  - Location と EVSE の1対多リレーションを含むレスポンス
  - CSVファイルからのサンプルデータインポート
  - Goによる距離計算ロジック（Haversine公式）

- **含まれないもの**:
  - 認証・認可機能
  - ロケーションの作成・更新・削除機能
  - フロントエンドUI
  - ページネーション（仕様書に記載なし）

### 1.3 新機能追加の詳細

#### 追加機能の動作

- **リクエスト受付**: `GET /api/locations` でクエリパラメータを受け取る
  - `latitude` (必須): 検索中心点の緯度（例: 50.770774）
  - `longitude` (必須): 検索中心点の経度（例: -126.104965）
  - `radius` (任意): 検索半径（km）、デフォルト100km
  - `date_from` (任意): この日時以降に更新されたロケーションのみ（inclusive）
  - `date_to` (任意): この日時より前に更新されたロケーションのみ（exclusive）

- **検索処理**:
  1. パラメータのバリデーション（正規表現チェック）
  2. データベースから全ロケーションを取得
  3. Goで各ロケーションとの距離を計算（Haversine公式）
  4. 指定半径内のロケーションをフィルタリング
  5. 日時フィルターの適用（該当する場合）
  6. EVSEを含むロケーション情報を返却

- **レスポンス形式**:
  ```json
  [
    {
      "id": "1",
      "name": "旭川日産自動車 枝幸店",
      "address": "北海道枝幸郡枝幸町南浜町1345-10",
      "coordinates": {
        "latitude": "44.922672",
        "longitude": "142.581987"
      },
      "evses": [
        {
          "uid": "CP010000048501",
          "status": "AVAILABLE"
        }
      ]
    }
  ]
  ```

#### 影響範囲

- **新規作成するファイル**:
  - `models/location.go` - Location モデル定義
  - `models/evse.go` - EVSE モデル定義
  - `handlers/location_handler.go` - HTTPハンドラー
  - `services/location_service.go` - ビジネスロジック
  - `database/seed.go` - CSVインポート処理

- **修正が必要な既存ファイル**:
  - `router/router.go` - エンドポイント登録
  - `server.go` - データベース初期化処理

#### 既存機能との関係

- **活用する既存実装パターン**:
  - Gin フレームワークによるルーティング
  - GORM によるデータベースアクセス
  - 既存のデータベース接続設定（config/database.go）

---

## 2. データモデル設計

### 2.1 Location テーブル

| カラム名 | データ型 | 制約 | 説明 |
|---------|---------|------|------|
| id | string(36) | PK | ロケーション固有ID |
| name | string(255) | NULL許可 | ロケーション名 |
| address | string(45) | NOT NULL | 住所 |
| latitude | float64 | NOT NULL | 緯度（数値として保存） |
| longitude | float64 | NOT NULL | 経度（数値として保存） |
| last_updated | datetime | NULL許可 | 最終更新日時 |

### 2.2 EVSE テーブル

| カラム名 | データ型 | 制約 | 説明 |
|---------|---------|------|------|
| uid | string(36) | PK | EVSE固有ID |
| location_id | string(36) | FK | 所属ロケーションID |
| status | int | NOT NULL | ステータス（1-9） |

### 2.3 Status 列挙値

| 値 | DB値 | 説明 |
|----|------|------|
| AVAILABLE | 1 | 利用可能 |
| BLOCKED | 2 | 物理的障害あり |
| CHARGING | 3 | 使用中 |
| INOPERATIVE | 4 | 一時的に利用不可 |
| OUTOFORDER | 5 | 故障中 |
| PLANNED | 6 | 計画中 |
| REMOVED | 7 | 撤去済み |
| RESERVED | 8 | 予約済み |
| UNKNOWN | 9 | 不明 |

---

## 3. API仕様

### 3.1 GET /api/locations

#### リクエストパラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|-----|------|------|
| latitude | string(11) | Yes | 緯度。正規表現: `-?[0-9]{1,2}\.[0-9]{5,7}` |
| longitude | string(12) | Yes | 経度。正規表現: `-?[0-9]{1,3}\.[0-9]{5,7}` |
| radius | int | No | 検索半径（km）。デフォルト: 100 |
| date_from | DateTime | No | この日時以降の更新のみ（inclusive） |
| date_to | DateTime | No | この日時より前の更新のみ（exclusive） |

#### レスポンス

- **200 OK**: Location配列（EVSEを含む）
- **400 Bad Request**: パラメータ不正

---

## 4. 技術的考慮事項

### 4.1 距離計算

Haversine公式を使用して2点間の距離を計算する。この計算はGoで実装し、データベースの地理空間拡張（PostGISなど）は使用しない（要件に「フィルタリングロジックは主にGoで実装」と記載）。

```go
// Haversine公式による距離計算（km）
func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
    const R = 6371 // 地球の半径（km）

    dLat := (lat2 - lat1) * math.Pi / 180
    dLon := (lon2 - lon1) * math.Pi / 180

    a := math.Sin(dLat/2)*math.Sin(dLat/2) +
        math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
        math.Sin(dLon/2)*math.Sin(dLon/2)

    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

    return R * c
}
```

### 4.2 パフォーマンス考慮

- サンプルデータは100件程度と小規模なため、全件取得後のGoでのフィルタリングで十分
- 大規模データの場合は、事前にバウンディングボックスでSQLフィルタリングを検討

### 4.3 テスタビリティ

- Handlerとサービス層を分離し、依存性注入を可能に
- インターフェースを使用してモック可能な設計

---

_この基本設計は詳細設計・実装の基盤となります。変更は適切な承認プロセスを経て行ってください。_
