package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/coco-papiyon/mechacopy/cmd"
	"github.com/coco-papiyon/mechacopy/worker"
)

func main() {
	usage := []string{
		fmt.Sprintf("使い方: %s [オプション] コピー元 コピー先 [ファイル [ファイル]...]", os.Args[0]),
		"           コピー元 :: コピー元ディレクトリ",
		"           コピー先 :: コピー先ディレクトリ",
		"           ファイル :: コピーするファイル (名前/ワイルドカード: 既定値は「*」",
	}

	// 引数を取得
	config, args := cmd.Args(2, strings.Join(usage, "\n"))
	src := args[0]
	dst := args[1]
	extraArgs := args[2:]

	// 動作設定
	if len(extraArgs) > 0 {
		config.TargetFiles = extraArgs
	}
	var runner worker.Runner = &worker.CopyRunner{
		Destination: dst,
	}

	// コピーを実行
	slog.Info("Start Copy", "コピー元", src, "コピー先", dst, "対象ファイル", config.TargetFiles)
	worker.RunMecha(src, runner, config)
}
