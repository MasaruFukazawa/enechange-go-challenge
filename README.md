ENECHANGE Go Developer Challenge
====

EV充電器ロケーション検索APIの実装

## 概要

指定された緯度・経度・半径で定義された円形エリア内のEV充電器ロケーションを検索するAPIです。

## 機能

- 緯度・経度・半径による円形範囲検索
- 更新日時によるフィルタリング（date_from, date_to）
- Location と EVSE の1対多リレーションを含むレスポンス
- Haversine公式による距離計算

## セットアップ

### 前提条件

- Docker
- Docker Compose

### 起動方法

```bash
# コンテナ起動
docker-compose up -d

# ログ確認
docker-compose logs -f app
```

サーバーは `http://localhost:8070` で起動します。

## API仕様

### GET /api/locations

EV充電器ロケーションを検索

#### リクエストパラメータ

| パラメータ | 型 | 必須 | 説明 |
|-----------|-----|------|------|
| latitude | string | Yes | 緯度（例: 35.676200） |
| longitude | string | Yes | 経度（例: 139.650300） |
| radius | int | No | 検索半径 km（デフォルト: 100） |
| date_from | datetime | No | この日時以降の更新のみ（RFC3339形式） |
| date_to | datetime | No | この日時より前の更新のみ（RFC3339形式） |

#### レスポンス例

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

#### 使用例

```bash
# 旭川付近50km圏内を検索
curl "http://localhost:8070/api/locations?latitude=43.796611&longitude=142.375917&radius=50"
```

## プロジェクト構成

```
├── main.go                 # エントリーポイント
├── server.go               # サーバー設定
├── config/
│   └── config.go           # 環境設定
├── database/
│   ├── database.go         # DB接続
│   └── seed.go             # CSVインポート
├── models/
│   ├── location.go         # Locationモデル
│   └── evse.go             # EVSEモデル
├── repositories/
│   └── location_repository.go  # データアクセス層
├── services/
│   ├── location_service.go # ビジネスロジック
│   └── distance.go         # 距離計算（Haversine）
├── handlers/
│   └── location_handler.go # HTTPハンドラー
├── router/
│   └── router.go           # ルーティング
├── e2e/
│   └── api_test.sh         # APIテストスクリプト
└── sample/
    ├── locations.csv       # サンプルデータ
    └── evses.csv           # サンプルデータ
```

## テスト

### ユニットテスト

```bash
# 全テスト実行
docker exec go-challenge-app go test ./... -v
```

### 統合テスト

```bash
# 統合テスト実行
docker exec go-challenge-app go test ./handlers/... -tags=integration -v
```

### APIテスト

```bash
# シェルスクリプトでAPIテスト（要: jq）
./e2e/api_test.sh

# curlで手動テスト
curl -s "http://localhost:8070/api/locations?latitude=43.796611&longitude=142.375917&radius=50" | python3 -c "import sys,json; d=json.load(sys.stdin); print(f'{len(d)} locations found')"
```

## 技術スタック

- **言語**: Go 1.23
- **Webフレームワーク**: Gin
- **ORM**: GORM
- **データベース**: MySQL 8.0

## 設計ドキュメント

- [基本設計](.specs/ev-charger-location-search/basic-design.md)
- [詳細設計](.specs/ev-charger-location-search/detailed-design.md)
- [OpenAPI仕様](openapi/openapi.yml)

---

## Challenge Requirements

### Objectives

Create an endpoint that meets the following requirements using the specified language, library and template.

Please document any technical decisions, trade-offs, problems, etc., in REPORT.md.

#### Requirements

- Specification: [Go-Challenge Interface Specification Document](./Go-Challenge%20Interface%20Specification%20Document.pdf)
- CSV Files:
  - [locations.csv](./sample/locations.csv)
  - [evses.csv](./sample/evses.csv)

#### Language / Libraries

- Language: Go
- Web Framework: [Gin](https://gin-gonic.com/)
- ORM: [GORM](https://gorm.io/)
