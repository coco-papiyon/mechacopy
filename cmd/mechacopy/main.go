package main

import (
	"log/slog"

	"github.com/coco-papiyon/mechacopy/cmd"
	"github.com/coco-papiyon/mechacopy/worker"
)

func main() {
	// 引数を取得
	config, args := cmd.Args(2)
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
