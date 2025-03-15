package worker

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/coco-papiyon/mechacopy/directory"
)

type Config struct {
	CopyThread  int
	RetryCount  int
	SleepTime   int
	TargetFiles []string
	Retry       bool
}

func InitConfig() *Config {
	return &Config{
		CopyThread:  10,
		RetryCount:  10,
		SleepTime:   10,
		TargetFiles: []string{"*"},
		Retry:       true,
	}
}

type Runner interface {
	Run(string, string, *JobStatus) error
	Retry(string, string) error
}

// ディレクトリ内のファイルをすべてコピーする(同時実行制御)
func RunMecha(srcDir string, runner Runner, config *Config) error {
	start := time.Now()

	// 指定ディレクトリ内のディレクトリを取得
	srcDirs, err := directory.GetDirs(srcDir)
	if err != nil {
		slog.Error("ディレクトリ取得", "ERROR", err, "basePath", srcDir)
		return err
	}

	// 同時実行用の制御
	job := &JobStatus{}
	job.totalCnt = int32(len(srcDirs))
	job.ch = make(chan string)
	job.wg.Add(len(srcDirs))
	job.config = config

	// 指定した数スレッド(goroutine)を起動
	for i := 0; i < config.CopyThread; i++ {
		go runWorker(srcDir, runner, job)
	}

	// コピー対象のディレクトリを送信
	for _, name := range srcDirs {
		job.ch <- name
	}

	// 処理待ち
	job.wg.Wait()

	// エラーリトライ
	if config.Retry {
		runRetry(srcDir, runner, job)
	}

	// 処理時間を取得
	end := time.Now()
	duration := end.Sub(start)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	errCnt := int32(len(job.errorFiles))
	slog.Info("All File Finished")
	fmt.Printf("処理時間: %02d時間 %02d分 %02d秒\n", hours, minutes, seconds)
	fmt.Printf("    Total:   %d\n", job.successFileCnt+job.skipFileCnt+errCnt)
	fmt.Printf("    Success: %d\n", job.successFileCnt)
	fmt.Printf("    Skip:    %d\n", job.skipFileCnt)
	fmt.Printf("    Error:   %d\n", errCnt)

	if errCnt > 0 {
		fmt.Printf("ERROR Files")
		for _, file := range job.errorFiles {
			fmt.Printf("  %s\n", file)
		}
	}
	return nil
}

// エラーとなったファイルのコピーをリトライ
func runRetry(srcDir string, runner Runner, job *JobStatus) {
	for i := 0; i < job.config.RetryCount; i++ {
		errCount := len(job.errorFiles)
		if errCount == 0 {
			break
		}

		slog.Info("Retry Error Files", "Count", i, "Files", errCount)
		time.Sleep(time.Duration(job.config.SleepTime) * time.Second)

		// エラーとなったファイルのコピーをリトライ
		errFiles := make([]string, errCount)
		copy(errFiles, job.errorFiles)
		job.errorFiles = []string{}
		job.wg = sync.WaitGroup{}
		job.wg.Add(len(errFiles))
		job.ch = make(chan string)

		// 指定した数スレッド(goroutine)を起動
		for i := 0; i < job.config.CopyThread; i++ {
			go retryWorker(srcDir, runner, job)
		}

		// コピー処理を実行
		for _, file := range errFiles {
			job.ch <- file
		}
		job.wg.Wait()
	}
}

// 非同期でコピーを行う
func retryWorker(srcDir string, runner Runner, job *JobStatus) {
	for {
		targetFile := <-job.ch
		err := runner.Retry(srcDir, targetFile)
		if err != nil {
			slog.Error("File Copy", "file", targetFile, "ERROR", err)
			job.AddErrorFile(targetFile)
		}
		job.wg.Done()
	}
}

// ファイルをコピーするワーカー(同時実行制御)
func runWorker(baseDir string, runner Runner, job *JobStatus) {
	for {
		srcDir := <-job.ch
		err := runner.Run(baseDir, srcDir, job)
		if err != nil {
			job.AddError()
			slog.Error(fmt.Sprintf("Copy Error: %s %s %v", job.GetStatus(), srcDir, err))
		} else {
			job.AddSuccess()
			slog.Info(fmt.Sprintf("%s %s", job.GetStatus(), srcDir))
		}
		job.wg.Done()
	}
}
