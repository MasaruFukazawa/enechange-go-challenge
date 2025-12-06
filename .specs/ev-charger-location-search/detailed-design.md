# 詳細設計: ev-charger-location-search

## 1. コンポーネント設計

### 1.1 データモデル層

#### M-001: Location モデル

**ファイル**: `models/location.go`

**責務**: ロケーションエンティティの定義とGORMマッピング

**構造体定義**:

```go
type Location struct {
    ID          string    `gorm:"primaryKey;size:36" json:"id"`
    Name        *string   `gorm:"size:255" json:"name"`
    Address     string    `gorm:"size:45;not null" json:"address"`
    Latitude    float64   `gorm:"not null" json:"-"`
    Longitude   float64   `gorm:"not null" json:"-"`
    LastUpdated *time.Time `json:"-"`
    EVSEs       []EVSE    `gorm:"foreignKey:LocationID" json:"evses"`
}
```

**JSON出力構造**:

```go
type LocationResponse struct {
    ID          string           `json:"id"`
    Name        *string          `json:"name,omitempty"`
    Address     string           `json:"address"`
    Coordinates CoordinatesDTO   `json:"coordinates"`
    EVSEs       []EVSEResponse   `json:"evses"`
}

type CoordinatesDTO struct {
    Latitude  string `json:"latitude"`
    Longitude string `json:"longitude"`
}
```

**メソッド**:

- `TableName() string` - テーブル名「locations」を返す
- `ToResponse() LocationResponse` - JSON出力用構造体への変換

---

#### M-002: EVSE モデル

**ファイル**: `models/evse.go`

**責務**: EVSEエンティティの定義とステータス列挙

**構造体定義**:

```go
type EVSE struct {
    UID        string `gorm:"primaryKey;size:36" json:"uid"`
    LocationID string `gorm:"size:36;not null" json:"-"`
    Status     int    `gorm:"not null" json:"-"`
}
```

**ステータス列挙**:

```go
type EVSEStatus int

const (
    StatusAvailable   EVSEStatus = 1
    StatusBlocked     EVSEStatus = 2
    StatusCharging    EVSEStatus = 3
    StatusInoperative EVSEStatus = 4
    StatusOutOfOrder  EVSEStatus = 5
    StatusPlanned     EVSEStatus = 6
    StatusRemoved     EVSEStatus = 7
    StatusReserved    EVSEStatus = 8
    StatusUnknown     EVSEStatus = 9
)

var statusStrings = map[EVSEStatus]string{
    StatusAvailable:   "AVAILABLE",
    StatusBlocked:     "BLOCKED",
    StatusCharging:    "CHARGING",
    StatusInoperative: "INOPERATIVE",
    StatusOutOfOrder:  "OUTOFORDER",
    StatusPlanned:     "PLANNED",
    StatusRemoved:     "REMOVED",
    StatusReserved:    "RESERVED",
    StatusUnknown:     "UNKNOWN",
}
```

**JSON出力構造**:

```go
type EVSEResponse struct {
    UID    string `json:"uid"`
    Status string `json:"status"`
}
```

**メソッド**:

- `TableName() string` - テーブル名「evses」を返す
- `(s EVSEStatus) String() string` - ステータスの文字列表現
- `ToResponse() EVSEResponse` - JSON出力用構造体への変換

---

### 1.2 リポジトリ層

#### R-001: LocationRepository

**ファイル**: `repositories/location_repository.go`

**責務**: Locationエンティティのデータアクセス

**インターフェース定義**:

```go
type LocationRepository interface {
    FindAll() ([]Location, error)
    FindAllWithEVSEs() ([]Location, error)
}
```

**実装**:

```go
type locationRepository struct {
    db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) LocationRepository {
    return &locationRepository{db: db}
}
```

**メソッド詳細**:

| メソッド | パラメータ | 戻り値 | 説明 |
|---------|-----------|--------|------|
| `FindAll` | なし | `([]Location, error)` | 全ロケーションを取得 |
| `FindAllWithEVSEs` | なし | `([]Location, error)` | EVSE含む全ロケーションを取得（Preload） |

---

### 1.3 サービス層

#### S-001: LocationService

**ファイル**: `services/location_service.go`

**責務**: ロケーション検索のビジネスロジック

**インターフェース定義**:

```go
type LocationService interface {
    SearchByRadius(params SearchParams) ([]LocationResponse, error)
}

type SearchParams struct {
    Latitude  float64
    Longitude float64
    Radius    float64      // km
    DateFrom  *time.Time   // inclusive
    DateTo    *time.Time   // exclusive
}
```

**依存関係**:

- `LocationRepository` - データアクセス

**実装クラス**:

```go
type locationService struct {
    repo LocationRepository
}

func NewLocationService(repo LocationRepository) LocationService {
    return &locationService{repo: repo}
}
```

**メソッド詳細**:

| メソッド | パラメータ | 戻り値 | 説明 |
|---------|-----------|--------|------|
| `SearchByRadius` | `SearchParams` | `([]LocationResponse, error)` | 半径検索＋日時フィルタリング |

**処理フロー**:

1. リポジトリから全ロケーション（EVSE含む）を取得
2. 各ロケーションに対して:
   - Haversine公式で距離を計算
   - 指定半径内かチェック
   - 日時フィルター適用（DateFrom/DateTo）
3. 条件を満たすロケーションをレスポンス形式に変換
4. 結果を返却

---

#### S-002: DistanceCalculator

**ファイル**: `services/distance.go`

**責務**: 地理的距離の計算

**関数定義**:

```go
// HaversineDistance は2点間の距離をkmで計算する
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
    const EarthRadiusKm = 6371.0

    lat1Rad := lat1 * math.Pi / 180
    lat2Rad := lat2 * math.Pi / 180
    deltaLat := (lat2 - lat1) * math.Pi / 180
    deltaLon := (lon2 - lon1) * math.Pi / 180

    a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
        math.Cos(lat1Rad)*math.Cos(lat2Rad)*
        math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

    return EarthRadiusKm * c
}
```

---

### 1.4 ハンドラー層

#### H-001: LocationHandler

**ファイル**: `handlers/location_handler.go`

**責務**: HTTPリクエストの処理とレスポンス生成

**構造体定義**:

```go
type LocationHandler struct {
    service LocationService
}

func NewLocationHandler(service LocationService) *LocationHandler {
    return &LocationHandler{service: service}
}
```

**依存関係**:

- `LocationService` - ビジネスロジック

**エンドポイント**:

| メソッド | パス | ハンドラーメソッド |
|---------|------|-------------------|
| GET | `/api/locations` | `GetLocations` |

**リクエストバリデーション**:

```go
type GetLocationsRequest struct {
    Latitude  string `form:"latitude" binding:"required"`
    Longitude string `form:"longitude" binding:"required"`
    Radius    int    `form:"radius"`
    DateFrom  string `form:"date_from"`
    DateTo    string `form:"date_to"`
}
```

**バリデーションルール**:

| パラメータ | 正規表現 | エラーメッセージ |
|-----------|---------|----------------|
| latitude | `-?[0-9]{1,2}\.[0-9]{5,7}` | Invalid latitude format |
| longitude | `-?[0-9]{1,3}\.[0-9]{5,7}` | Invalid longitude format |
| radius | 正の整数 | Invalid radius value |
| date_from | RFC3339形式 | Invalid date_from format |
| date_to | RFC3339形式 | Invalid date_to format |

**エラーレスポンス**:

```go
type ErrorResponse struct {
    Error string `json:"error"`
}
```

**処理フロー**:

1. クエリパラメータのバインディング
2. 必須パラメータの存在チェック
3. 正規表現によるフォーマットバリデーション
4. 日時パラメータのパース（存在する場合）
5. SearchParamsの構築
6. サービス呼び出し
7. レスポンス返却（200 OK または 400 Bad Request）

---

### 1.5 データベース初期化

#### DB-001: Seeder

**ファイル**: `database/seed.go`

**責務**: CSVファイルからの初期データ投入

**関数定義**:

```go
func SeedFromCSV(db *gorm.DB, locationsPath, evsesPath string) error
```

**処理フロー**:

1. テーブルの存在確認・作成（AutoMigrate）
2. 既存データの削除（冪等性確保）
3. locations.csvの読み込み・パース
4. Locationレコードの一括挿入
5. evses.csvの読み込み・パース
6. EVSEレコードの一括挿入

**CSVフォーマット**:

locations.csv:
```
id,name,address,latitude,longitude,last_updated
```

evses.csv:
```
uid,location_id,status
```

---

### 1.6 ルーター設定

#### RT-001: Router

**ファイル**: `router/router.go`

**修正内容**:

```go
func SetupRouter(db *gorm.DB) *gin.Engine {
    r := gin.Default()

    // 依存性の注入
    locationRepo := repositories.NewLocationRepository(db)
    locationService := services.NewLocationService(locationRepo)
    locationHandler := handlers.NewLocationHandler(locationService)

    // エンドポイント登録
    api := r.Group("/api")
    {
        api.GET("/locations", locationHandler.GetLocations)
    }

    return r
}
```

---

## 2. シーケンス図

### 2.1 ロケーション検索フロー

```
Client          Handler           Service          Repository        DB
  |                |                 |                 |              |
  |-- GET /api/locations ---------->|                 |              |
  |                |                 |                 |              |
  |                |-- validate params                |              |
  |                |                 |                 |              |
  |   400 Bad Request (if invalid)  |                 |              |
  |<---------------|                 |                 |              |
  |                |                 |                 |              |
  |                |-- SearchByRadius -------------->|              |
  |                |                 |                 |              |
  |                |                 |-- FindAllWithEVSEs ---------->|
  |                |                 |                 |              |
  |                |                 |<-- []Location ---------------+|
  |                |                 |                 |              |
  |                |                 |-- filter by distance          |
  |                |                 |-- filter by date              |
  |                |                 |-- convert to response         |
  |                |                 |                 |              |
  |                |<-- []LocationResponse            |              |
  |                |                 |                 |              |
  |<-- 200 OK -----|                 |                 |              |
  |   (JSON array) |                 |                 |              |
```

---

## 3. ファイル構成

```
enechange-go-challenge/
├── main.go                          # エントリーポイント
├── server.go                        # サーバー設定（修正）
├── config/
│   └── database.go                  # DB接続設定（既存）
├── models/
│   ├── location.go                  # 新規作成
│   └── evse.go                      # 新規作成
├── repositories/
│   └── location_repository.go       # 新規作成
├── services/
│   ├── location_service.go          # 新規作成
│   └── distance.go                  # 新規作成
├── handlers/
│   └── location_handler.go          # 新規作成
├── database/
│   └── seed.go                      # 新規作成
├── router/
│   └── router.go                    # 修正
└── data/
    ├── locations.csv                # サンプルデータ（既存）
    └── evses.csv                    # サンプルデータ（既存）
```

---

## 4. エラーハンドリング

### 4.1 HTTPエラーコード

| コード | 条件 | レスポンス例 |
|-------|------|-------------|
| 200 | 正常終了 | `[{...}, {...}]` |
| 400 | パラメータ不正 | `{"error": "Invalid latitude format"}` |
| 500 | サーバーエラー | `{"error": "Internal server error"}` |

### 4.2 バリデーションエラー詳細

| エラー種別 | 検出タイミング | メッセージ |
|-----------|--------------|-----------|
| 必須パラメータ欠落 | Handler | `latitude is required` / `longitude is required` |
| フォーマット不正 | Handler | `Invalid latitude format` / `Invalid longitude format` |
| 日時パース失敗 | Handler | `Invalid date_from format` / `Invalid date_to format` |
| 半径不正値 | Handler | `Invalid radius value` |

---

## 5. テスト設計

### 5.1 ユニットテスト

#### 距離計算テスト (`services/distance_test.go`)

| テストケース | 入力 | 期待値 |
|-------------|------|--------|
| 同一地点 | (35.6762, 139.6503) - (35.6762, 139.6503) | 0 km |
| 東京-大阪 | (35.6762, 139.6503) - (34.6937, 135.5023) | 約395 km |
| 北半球-南半球 | (35.6762, 139.6503) - (-33.8688, 151.2093) | 約7823 km |

#### サービステスト (`services/location_service_test.go`)

| テストケース | 条件 | 期待結果 |
|-------------|------|---------|
| 半径内のみ取得 | 半径100km | 範囲内ロケーションのみ |
| 空の結果 | 半径1km、近くにロケーションなし | 空配列 |
| 日時フィルター | date_from指定 | 指定日時以降のみ |

#### ハンドラーテスト (`handlers/location_handler_test.go`)

| テストケース | リクエスト | 期待ステータス |
|-------------|-----------|---------------|
| 正常リクエスト | latitude=35.6762&longitude=139.6503 | 200 |
| latitude欠落 | longitude=139.6503 | 400 |
| 不正フォーマット | latitude=invalid | 400 |

### 5.2 統合テスト

| テストケース | 説明 | 検証内容 |
|-------------|------|---------|
| E2E検索 | CSVデータからの検索 | レスポンスJSONの形式とデータ |
| EVSE含有確認 | ロケーションにEVSE含む | evses配列が正しく含まれる |

---

## 6. 技術的リスクと対策

### 6.1 パフォーマンス

| リスク | 影響度 | 対策 |
|-------|-------|------|
| 全件取得の遅延 | 低（100件程度） | 現状で問題なし |
| 大量データ時の遅延 | 中 | バウンディングボックスによる事前フィルタリング検討 |

### 6.2 精度

| リスク | 影響度 | 対策 |
|-------|-------|------|
| Haversine公式の誤差 | 低 | 楕円体モデル（Vincenty）は不要と判断 |
| 浮動小数点誤差 | 低 | float64で十分な精度 |

---

_この詳細設計は基本設計に基づいて作成されています。実装時はこのドキュメントを参照してください。_
