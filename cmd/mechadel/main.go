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
		fmt.Sprintf("使い方: %s [オプション] 削除対象", os.Args[0]),
		"           削除対象 :: 削除対象のディレクトリ",
	}

	// 引数を取得
	config, args := cmd.Args(1, strings.Join(usage, "\n"))
	src := args[0]

	// 動作設定
	config.Retry = false
	var runner worker.Runner = &worker.DeleteRunner{}

	slog.Info("Start Delete", "削除対象", src)
	worker.RunMecha(src, runner, config)
}
