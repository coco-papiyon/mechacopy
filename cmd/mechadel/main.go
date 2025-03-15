package main

import (
	"log/slog"

	"github.com/coco-papiyon/mechacopy/cmd"
	"github.com/coco-papiyon/mechacopy/worker"
)

func main() {
	// 引数を取得
	config, args := cmd.Args(1)
	src := args[0]

	// 動作設定
	config.Retry = false
	var runner worker.Runner = &worker.DeleteRunner{}

	slog.Info("Start Delete", "削除対象", src)
	worker.RunMecha(src, runner, config)
}
