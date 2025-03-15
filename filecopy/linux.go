//go:build !windows

package filecopy

import (
	"os"
)

// 作成日時、更新日時をコピーする
func copyTimestamps(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	modTime := srcInfo.ModTime()
	accessTime := modTime

	// Linux/macOS: os.Chtimes() を使用
	return os.Chtimes(dst, accessTime, modTime)
}
