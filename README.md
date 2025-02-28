# シンプルな Go プロジェクト

これはシンプルな Go プロジェクトのテンプレートです。

## 概要

このリポジトリは基本的な Go プロジェクトの構造を示しています。

## インストール方法

```bash
# リポジトリをクローン
git clone https://github.com/yourusername/simple-go-project.git

# ディレクトリに移動
cd simple-go-project

# 依存関係をインストール（必要な場合）
go mod download
```

## 使用方法

```bash
# アプリケーションを実行
go run ./cmd/main.go
```

## プロジェクト構造

```
.
├── cmd/            # メインアプリケーションのエントリーポイント
├── internal/       # 外部からインポートされるべきでないパッケージ
├── pkg/            # 外部からインポートできるパッケージ
├── go.mod          # Goモジュール定義
├── go.sum          # 依存関係のチェックサム
└── README.md       # このファイル
```

## ライセンス

MIT
