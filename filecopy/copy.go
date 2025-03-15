package filecopy

import (
	"io"
	"os"
	"path/filepath"
)

// Config.TargetFilesに一致するファイルかチェックする
func IsCopyFile(srcFile string, targetFiles []string) bool {
	// ファイル名のチェック対象
	srcFile = filepath.Base(srcFile)
	for _, pattern := range targetFiles {
		matched, _ := filepath.Match(pattern, srcFile)
		//fmt.Println(pattern, srcFile, matched)
		if matched {
			return true
		}
	}
	return false
}

// ファイルの差分をチェックする(サイズ、更新日付)
func IsFileDiff(src, dst string) bool {
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

// ファイルコピーし、更新日時等変更する
func CopyFile(src, dst string) error {
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
