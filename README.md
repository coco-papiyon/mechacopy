# Mechacopy

Mechacopyは、Go言語で実装されたファイルコピー・削除・ディレクトリ操作などを行うツールです。

## プロジェクト構成

- `cmd/` : コマンドラインツールのエントリーポイント
  - `mechacopy/main.go` : メインコピーコマンド
  - `mechadel/main.go` : 削除コマンド（mechadelコマンドは指定したファイルやディレクトリを安全に削除するためのツールです。コマンドライン引数で削除対象を指定できます）
- `directory/` : ディレクトリ操作関連
- `filecopy/` : ファイルコピー関連
- `testutil/` : テスト用ユーティリティ
- `worker/` : ワーカー処理（コピー・削除ジョブ管理）

## 動作環境
- Windows  
- Linux  

## 使い方(コピー)  

### ビルド
```sh
go build -o mechacopy ./cmd/mechacopy
```

### 実行例
```sh
./mechacopy <コピー元> <コピー先> <対象ファイル名>
```

## 使い方(削除)  

### ビルド
```sh
go build -o mechadel ./cmd/mechadel
```

### 実行例
```sh
./mechadel <コピー削除対象> 
```

## テスト
```sh
go test ./...
```

## ライセンス
MIT
