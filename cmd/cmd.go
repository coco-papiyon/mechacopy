package cmd

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/coco-papiyon/mechacopy/worker"
)

func Args(argLen int, usage string) (*worker.Config, []string) {
	var logger *slog.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// オプションの定義
	mt := flag.Int("MT", 0, "n 個のスレッドのマルチスレッド コピーを実行する (既定値 10)")
	retry := flag.Int("R", 0, "失敗したコピーに対する再試行数 (既定値 10)")
	wait := flag.Int("W", 0, "試行と再試行の間の待機時間 (既定値 10)")

	// Usageの出力
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n\n", usage)
		flag.PrintDefaults()
	}

	// オプションを解析
	flag.Parse()

	args := flag.Args()
	// 必須の引数をチェック（必須は2）
	if len(args) < argLen {
		flag.Usage()
		os.Exit(1)
	}

	config := worker.InitConfig()
	if *mt > 0 {
		config.CopyThread = *mt
		logger.Info(fmt.Sprintf("スレッド数: %d", *mt))
	}
	if *retry > 0 {
		config.RetryCount = *retry
		logger.Info(fmt.Sprintf("リトライ回数: %d", *retry))
	}
	if *wait > 0 {
		config.SleepTime = *wait
		logger.Info(fmt.Sprintf("リトライ待機時間: %d", *wait))
	}
	return config, args
}
