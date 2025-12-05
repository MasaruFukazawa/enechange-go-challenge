# 詳細設計: ev-charger-location-search

## コンポーネント設計

### UI コンポーネント

（本機能はAPIのみのため、UIコンポーネントなし）

### ビジネスロジックコンポーネント

#### BL-001: LocationHandler

**責務**: HTTPリクエストの受付とレスポンス返却

**依存関係**:

- LocationService

**公開メソッド**:

- GetLocations(c *gin.Context) - ロケーション検索エンドポイント

**エラーハンドリング**:

- 400: パラメータバリデーションエラー
- 500: 内部サーバーエラー

#### BL-002: LocationService

**責務**: ロケーション検索のビジネスロジック

**依存関係**:

- LocationRepository

**公開メソッド**:

- SearchLocations(params SearchParams) ([]Location, error) - 検索条件に基づくロケーション取得

## 4. 技術リスクの列挙

- **技術的課題**: 距離計算のパフォーマンス（Haversine公式の実装）
- **依存関係リスク**: 大量データ時のレスポンス時間
