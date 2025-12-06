---
description: 機能/不具合の詳細から仕様を初期化し、包括的な基本設計を生成する
---

**Notion URL or Github repository URL**: $ARGUMENTS[0](required)
**Figma URL**: $ARGUMENTS[1](nullable)

### フェーズ１： 必要なディレクトリ・ファイルの生成

#### 1.1 機能名生成

Notion/Github(+Figma)に記載されている機能/不具合の詳細に基づいて、実装の本質を捉える簡潔で説明的な機能名を作成。

**機能名生成品質基準**:

- 簡潔性: 機能の本質を 3-5 語程度で表現
- 重複回避: 既存仕様名との重複確認

#### 1.2 仕様書のディレクトリ/ファイルを作成

`.specs/{生成された機能名}/`ディレクトリを作成:

- `basic-design.md` - 基本設計書（完全版）
- `detailed-design.md` - 詳細設計用テンプレート
- `tasks.md` - 実装タスク用テンプレート
- `spec.json` - メタデータと承認追跡

#### 1.3 spec.json メタデータ初期化

`.claude/templates/spec.json`の形式に従って、spec.json を初期化。

### フェーズ 2: 基本設計生成

`.claude/templates/basic-design.md`の形式に従って、basic-design.md を生成。

## 次のステップ

1. basic-design.md の承認完了を確認
2. `specs detailed-design {機能名}` を実行して詳細設計を生成
3. 詳細設計レビューと承認

---
