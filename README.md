# インターフェースサイズと凝集度の比較

このプロジェクトはインターフェースの大きさによるコードの凝集度の違いを証明するためのサンプルです。

## 概要

このリポジトリでは、以下の 2 つのアプローチを比較しています：

1. **Big Interface アプローチ**：

   - 大きな単一インターフェース（`datastore`）を使用
   - すべてのデータ操作を 1 つのインターフェースにまとめる
   - 凝集度が低く、不要な関数の Behavior も定義する必要がある

2. **Small Interface アプローチ**：
   - 小さく特化したインターフェース（`userstore`, `todostore`）を使用
   - 各ドメインロジックに必要な操作のみを定義
   - 高い凝集度で、必要なものだけを利用できる

## インストール方法

```bash
# リポジトリをクローン
git clone https://github.com/TakumaKurosawa/big-interface-vs-small-interface.git

# ディレクトリに移動
cd big-interface-vs-small-interface

# 依存関係をインストール
go mod download
```

## テスト実行

```bash
# すべてのテストを実行
go test ./...
```

## プロジェクト構造

```
.
├── cmd/                        # メインアプリケーションのエントリーポイント
├── internal/                   # 外部からインポートされるべきでないパッケージ
│   ├── domain/                 # ドメインモデル
│   ├── biginterface/           # Big Interfaceアプローチの定義
│   │   ├── datastore.go        # 大きなインターフェース定義
│   │   └── mocks/              # Big Interfaceのモック
│   │       └── mock_datastore.go # DataStoreモック
│   ├── smallinterface/         # Small Interfaceアプローチの定義
│   │   ├── userstore.go        # ユーザー関連の小さなインターフェース
│   │   ├── todostore.go        # Todo関連の小さなインターフェース
│   │   └── mocks/              # Small Interfaceのモック
│   │       ├── mock_userstore.go # UserStoreモック
│   │       └── mock_todostore.go # TodoStoreモック
│   └── services/               # サービス実装
│       ├── biginterface/       # Big Interfaceアプローチのサービス実装
│       │   ├── service.go      # サービス実装
│       │   └── service_test.go # サービステスト
│       └── smallinterface/     # Small Interfaceアプローチのサービス実装
│           ├── service.go      # サービス実装
│           └── service_test.go # サービステスト
├── pkg/                        # 外部からインポートできるパッケージ
│   └── greeting/               # 共通モジュール
├── go.mod                      # Goモジュール定義
├── go.sum                      # 依存関係のチェックサム
└── README.md                   # このファイル
```

## ライセンス

MIT
