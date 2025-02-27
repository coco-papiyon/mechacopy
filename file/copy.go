package file

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type CopyConfig struct {
	CopyThread  int
	RetryCount  int
	SleepTime   int
	TargetFiles []string
}

var Config = &CopyConfig{
	CopyThread:  10,
	RetryCount:  10,
	SleepTime:   10,
	TargetFiles: []string{"*"},
}

// ディレクトリ内のファイルをすべてコピーする(同時実行制御)
func CopyFiles(srcDir string, dstDir string) error {
	start := time.Now()

	// 指定ディレクトリ内のディレクトリを取得
	srcDirs, err := GetDirs(srcDir)
	if err != nil {
		slog.Error("ディレクトリ取得", "ERROR", err, "basePath", srcDir)
		return err
	}

	// 同時実行用の制御
	job := &jobStatus{}
	job.totalCnt = int32(len(srcDirs))
	job.ch = make(chan string)
	job.wg.Add(len(srcDirs))

	// 指定した数スレッド(goroutine)を起動
	for i := 0; i < Config.CopyThread; i++ {
		go copyFileWorker(srcDir, dstDir, job)
	}

	// コピー対象のディレクトリを送信
	for _, name := range srcDirs {
		job.ch <- name
	}

	// 処理待ち
	job.wg.Wait()

	// エラーリトライ
	fmt.Println(job.errorFiles)
	copyFileTry(srcDir, dstDir, job)

	// 処理時間を取得
	end := time.Now()
	duration := end.Sub(start)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	errCnt := int32(len(job.errorFiles))
	slog.Info("All File COPY Finished")
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
func copyFileTry(srcDir, dstDir string, job *jobStatus) {
	for i := 0; i < Config.RetryCount; i++ {
		errCount := len(job.errorFiles)
		if errCount == 0 {
			break
		}

		slog.Info("Retry Error Files", "Count", i, "Files", errCount)
		time.Sleep(time.Duration(Config.SleepTime) * time.Second)

		// エラーとなったファイルのコピーをリトライ
		errFiles := make([]string, errCount)
		copy(errFiles, job.errorFiles)
		job.errorFiles = []string{}
		job.wg = sync.WaitGroup{}
		job.wg.Add(len(errFiles))
		job.ch = make(chan string)

		// 指定した数スレッド(goroutine)を起動
		for i := 0; i < Config.CopyThread; i++ {
			go retryWorker(srcDir, dstDir, job)
		}

		// コピー処理を実行
		for _, file := range errFiles {
			job.ch <- file
		}
		job.wg.Wait()
	}
}

// 非同期でコピーを行う
func retryWorker(srcDir, dstDir string, job *jobStatus) {
	for {
		file := <-job.ch
		src := filepath.Join(srcDir, file)
		dst := filepath.Join(dstDir, file)
		err := copyFile(src, dst)
		if err != nil {
			slog.Error("File Copy", "file", src, "ERROR", err)
			job.addErrorFile(file)
		}
		job.wg.Done()
	}
}

// ファイルをコピーするワーカー(同時実行制御)
func copyFileWorker(baseDir, dstBaseDir string, job *jobStatus) {
	for {
		srcDir := <-job.ch
		dstDir := filepath.Join(dstBaseDir, srcDir)
		err := copyFiles(baseDir, srcDir, dstDir, job)
		if err != nil {
			job.addError()
			slog.Error(fmt.Sprintf("Copy Error: %s %s %v", job.getStatus(), srcDir, err))
		} else {
			job.addSuccess()
			slog.Info(fmt.Sprintf("%s %s", job.getStatus(), srcDir))
		}
		job.wg.Done()
	}
}

// Config.TargetFilesに一致するファイルかチェックする
func isCopyFile(srcFile string) bool {
	// ファイル名のチェック対象
	srcFile = filepath.Base(srcFile)
	for _, pattern := range Config.TargetFiles {
		matched, _ := filepath.Match(pattern, srcFile)
		//fmt.Println(pattern, srcFile, matched)
		if matched {
			return true
		}
	}
	return false
}

// ファイルの差分をチェックする(サイズ、更新日付)
func isFileDiff(src, dst string) bool {
	// 元ファイルの情報取得
	srcInfo, err := os.Stat(src)
	if err != nil {
		return true
	}

	// コピー先のファイルの情報取得
	dstInfo, err := os.Stat(dst)
	if err != nil {
		return true
	}

	// ファイルのサイズと更新日を比較
	if dstInfo.Size() == srcInfo.Size() {
		if dstInfo.ModTime().UnixNano() >= srcInfo.ModTime().UnixNano() {
			return false
		}
	}

	return true
}

// ディレクトリ内のファイルをコピーする関数（サブディレクトリは無視）
func copyFiles(baseDir, srcBaseDir, dstDir string, job *jobStatus) error {
	srcDir := filepath.Join(baseDir, srcBaseDir)

	// ディレクトリ内のファイル一覧を取得
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		job.addErrorDirs(srcDir)
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			// ファイルのみをコピーする
			srcFile := filepath.Join(srcDir, entry.Name())
			dstFile := filepath.Join(dstDir, entry.Name())

			// コピー対象のファイルではない場合はスキップする
			if !isCopyFile(srcFile) {
				job.addSkipFile()
				continue
			}

			// 差分がない場合はコピーしない(サイズ、更新日付
			diff := isFileDiff(srcFile, dstFile)
			if !diff {
				job.addSkipFile()
				continue
			}

			// 元ファイルの情報取得
			var size int64 = 0
			srcInfo, err := os.Stat(srcFile)
			if err == nil {
				if srcInfo.Size() > 1073741824 {
					size = srcInfo.Size() / 1073741824
				}
			}

			// ファイルサイズが大きい場合はログ出力
			if size > 0 {
				slog.Info("START COPY BIG FILE", "file", srcFile, "size(GB)", size)
			}

			// ファイルコピー
			err = copyFile(srcFile, dstFile)
			if err != nil {
				slog.Error("File Copy", "file", srcFile, "ERROR", err)
				errFile := filepath.Join(srcBaseDir, entry.Name())
				job.addErrorFile(errFile)
				continue
			}

			// ファイルサイズが大きい場合はログ出力
			if size > 0 {
				slog.Info("END COPY BIG FILE", "file", srcFile, "size(GB)", size)
			}

			job.addSuccessFile()
		}
	}

	return nil
}

// ファイルコピーし、更新日時等変更する
func copyFile(src, dst string) error {
	err := copyData(src, dst)
	if err != nil {
		return err
	}
	return copyTimestamps(src, dst)
}

// ファイルコピー
func copyData(src, dst string) error {
	// 出力先ディレクトリを作成
	dstDir := filepath.Dir(dst)
	err := os.MkdirAll(dstDir, os.ModePerm)
	if err != nil {
		return err
	}

	// 元ファイルを開く
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// コピー先ファイルを作成
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// データをコピー
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// ファイルのバッファをフラッシュ
	return dstFile.Sync()
}
