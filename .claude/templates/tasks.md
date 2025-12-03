# 実装タスク: XXX

## タスク概要

---

## コアタスク (CORE)

### CORE-001: XXX

- **説明**: XXX
- **受入条件**:
  - XXX
  - XXX
- **依存関係**: XXX
- **見積時間**: X 時間
- **ファイル**:
  - `path/to/file.dart`

---

## UI タスク (UI)

### UI-001: XXX

- **説明**: XXX
- **受入条件**:
  - XXX
  - XXX
  - XXX
  - XXX
- **依存関係**: UI-XXX
- **見積時間**: X 時間
- **ファイル**:
  - `path/to/widget.dart`

### UI-002: XXX

- **説明**: XXX
- **受入条件**:
  - XXX
  - XXX
  - XXX
- **依存関係**: UI-005
- **見積時間**: X 時間
- **ファイル**:
  - `path/to/page.dart`

### UI-003: XXX

- **説明**: XXX
- **受入条件**:
  - XXX
  - XXX
  - XXX
  - XXX
- **依存関係**: CORE-XXX
- **見積時間**: X 時間
- **ファイル**:
  - `path/to/file.dart` (削除)
  - `path/to/router.dart`

---

## ドキュメントタスク (DOC)

### DOC-001: XXX

- **説明**: XXX
- **受入条件**:
  - XXX
  - XXX
  - XXX
  - XXX
- **依存関係**: UI-XXX, ANALYTICS-XXX
- **見積時間**: X 時間
- **ファイル**:
  - `docs/features/XXX.md`

---

## タスク依存関係図

```
【フェーズ1: XXX】
CORE-001 (XXX)
    ↓
UI-003 (XXX)

【フェーズ2: XXX】
UI-003 → UI-001, UI-002 (XXX)

【フェーズ3: XXX】
UI-001, UI-002 → DOC-001
```
