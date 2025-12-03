---
description: $ARGUMENTSのタスクリストを生成する
---

## タスクリストの生成

@.specs/$ARGUMENTS/detailed-design.md を参考にして、`.claude/templates/tasks.md`の形式に従い、tasks.md を生成。
タスク構造, タスクカテゴリ, 注意点は以下を参照：

### タスク構造

- **ID**: 一意識別子（TASK-001、TASK-002 など）
- **タイトル**: タスクの概要
- **説明**: 実行すべき内容
- **依存関係**: 先に完了すべきタスク
- **ファイル**: 作成/修正予定のファイル

### タスクカテゴリ

1. **コア** (CORE-XXX)

   - データモデル
   - ビジネスロジック
   - コアサービス
   - ユーティリティ

2. **UI** (UI-XXX)

   - コンポーネント
   - ページ/画面
   - スタイリング
   - ナビゲーション

3. **ドキュメント** (DOC-XXX)
   - API ドキュメント
   - ユーザーガイド
   - 開発者ドキュメント
   - README 更新

### 注意点

- タスク ID を使用して依存関係を指定
- クリティカルパスタスクを強調

### テストケース列挙の準備

タスク生成後、後続のワークフローで必要なディレクトリとファイルを作成します。

```
.specs/$ARGUMENTS/
├── basic-design.md          (既存)
├── detailed-design.md       (既存)
├── spec.json               (既存)
├── tasks.md                (既存)
├── testcases/              (新規作成ディレクトリ)
└── task-progress.json      (新規作成ファイル)
```

また、`.claude/templates/task-progress.json`の形式に従って、task-progress.json を初期化。

### メタデータ更新

タスク生成とファイル初期化完了後の更新をしてください。

**注意**: タスクの詳細管理（進捗、見積時間、完了状況など）は全て task-progress.json で行います。spec.json は仕様レベルの承認状況のみを管理します。
