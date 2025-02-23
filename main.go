package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/coco-papiyon/mechacopy/file"
)

func main() {
	var logger *slog.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// オプションの定義
	mt := flag.Int("MT", 0, "n 個のスレッドのマルチスレッド コピーを実行する (既定値 10)")

	// Usageの出力
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "使い方: %s [オプション] コピー元 コピー先 [ファイル [ファイル]...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "           コピー元 :: コピー元ディレクトリ\n")
		fmt.Fprintf(os.Stderr, "           コピー先 :: コピー先ディレクトリ\n")
		fmt.Fprintf(os.Stderr, "           ファイル :: コピーするファイル (名前/ワイルドカード: 既定値は「*」\n")
		flag.PrintDefaults()
	}

	// オプションを解析
	flag.Parse()
	args := flag.Args()

	// 必須の引数をチェック（必須は2）
	if len(args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	// 引数を取得
	src := args[0]
	dst := args[1]
	extraArgs := args[2:]

	// コピー動作設定
	if len(extraArgs) > 0 {
		file.Config.TargetFiles = extraArgs
	}
	logger.Info("Start Copy", "コピー元", src, "コピー先", dst, "対象ファイル", file.Config.TargetFiles)

	if *mt > 0 {
		file.Config.CopyThread = *mt
		logger.Info("          ", "スレッド数", *mt)
	}
	file.CopyFiles(src, dst)
}
