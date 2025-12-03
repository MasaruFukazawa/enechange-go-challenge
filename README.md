ENECHANGE Go Developer Challenge
====

### Why do we ask you to do this challenge?

As part of the ENECHANGE interview process, we are asking you to work on a small coding exercise to help us understand your skills, and give you an idea of the work you would be doing with us.

### Objectives

Create an endpoint that meets the following requirements using the specified language, library and template.

Please document any technical decisions, trade-offs, problems, etc., in REPORT.md (you may write in either Japanese or English).

#### Requirements:
Please create an endpoint to search for chargers within a specified range, according to the interface specifications outlined in the provided PDF document.

However, latitude and longitude should be stored separately as numerical values and the filtering logic should primarily be implemented in Go.

You can use the provided CSV files as sample data by importing them into the database.

While writing tests is not mandatory, please be mindful to design for testability.

- Specification:
   - [Go-Challenge Interface Specification Document](./Go-Challenge%20Interface%20Specification%20Document.pdf)

- CSV Files:
    - [locations.csv](./sample/locations.csv)
    - [evses.csv](./sample/evses.csv)

#### Language / Libraries:
- Language: GO
- Libraries:
  - Web Application Framework: [Gin](https://gin-gonic.com/)
  - ORM: [GORM](https://gorm.io/)

#### Template:
A template has been prepared in the codebase that sets up a server and prepares the connection to a database.

Please add code to implement the endpoint.

#### Steps
1. Please fork this repository.
2. Create a new branch. (The branch name should be `challenge/YOURACCOUNT` where `YOURACCOUNT` is your Github Account name, for example, `challenge/shirakia`)
3. Proceed with your implementation using provided template.
4. After completion, create a Pull Request and please provide us with the URL.
5. You may continue to commit to your branch even after you have provided us with the URL of the Pull Request.

# Claude Code Knowledge

このプロジェクトは、開発に使える Claude Code のコマンドを集積させることを目的としたリポジトリです。各種開発フローやコマンドの整備、そしてそれらの使い方を記したドキュメントを体系的に整理・提供しています。

## ドキュメント

`docs/` ディレクトリには、各種開発フローの詳細な解説ドキュメントが格納されています。

### 利用可能なコマンド

- **[大・中規模タスク開発フロー](docs/task-large.md)** (`task-large/`) - TDD（テスト駆動開発）と人間承認プロセスを組み合わせた高品質なソフトウェア開発手法。新機能開発など大規模・中規模タスクの実行に使用され、仕様管理から実装まで体系化されたワークフローを提供します。

- **[小規模タスク開発フロー](docs/task-small.md)** (`task-small/`) - 小規模タスク向けの簡潔な AI 開発コマンド集。バグ修正や機能の小規模改修など、数時間～ 1 日程度のタスクに最適化された迅速で柔軟なワークフローを提供します。

- **[その他のコマンド](docs/etc.md)** (`etc/`) - 単体で動作する特化型コマンド群。リサーチ・分析（research、error-analysis）など、特定の用途に最適化されたツールを提供します。
