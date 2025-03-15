package worker

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/coco-papiyon/mechacopy/filecopy"
)

type CopyRunner struct {
	Destination string
}

func (r CopyRunner) Run(baseDir, srcDir string, job *JobStatus) error {
	dstDir := filepath.Join(r.Destination, srcDir)
	return copyFiles(baseDir, srcDir, dstDir, job)
}

func (r CopyRunner) Retry(srcDir, targetFile string) error {
	src := filepath.Join(srcDir, targetFile)
	dst := filepath.Join(r.Destination, targetFile)
	return filecopy.CopyFile(src, dst)
}

// ディレクトリ内のファイルをコピーする関数（サブディレクトリは無視）
func copyFiles(baseDir, srcBaseDir, dstDir string, job *JobStatus) error {
	srcDir := filepath.Join(baseDir, srcBaseDir)

	// ディレクトリ内のファイル一覧を取得
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		job.AddErrorDirs(srcDir)
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			// ファイルのみをコピーする
			srcFile := filepath.Join(srcDir, entry.Name())
			dstFile := filepath.Join(dstDir, entry.Name())

			// コピー対象のファイルではない場合はスキップする
			if !filecopy.IsCopyFile(srcFile, job.config.TargetFiles) {
				job.AddSkipFile()
				continue
			}

			// 差分がない場合はコピーしない(サイズ、更新日付
			diff := filecopy.IsFileDiff(srcFile, dstFile)
			if !diff {
				job.AddSkipFile()
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
			err = filecopy.CopyFile(srcFile, dstFile)
			if err != nil {
				slog.Error("File Copy", "file", srcFile, "ERROR", err)
				errFile := filepath.Join(srcBaseDir, entry.Name())
				job.AddErrorFile(errFile)
				continue
			}

			// ファイルサイズが大きい場合はログ出力
			if size > 0 {
				slog.Info("END COPY BIG FILE", "file", srcFile, "size(GB)", size)
			}

			job.AddSuccessFile()
		}
	}

	return nil
}
