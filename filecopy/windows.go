//go:build windows

package filecopy

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// Windows で作成日・更新日・アクセス日時をコピー
func copyTimestamps(src, dst string) error {
	// Windows の API を使ってタイムスタンプを取得
	srcPath, err := windows.UTF16PtrFromString(src)
	if err != nil {
		return err
	}

	var data windows.Win32FileAttributeData
	err = windows.GetFileAttributesEx(srcPath, windows.GetFileExInfoStandard, (*byte)(unsafe.Pointer(&data)))
	if err != nil {
		return err
	}

	// コピー先ファイルのハンドルを取得
	dstPath, err := windows.UTF16PtrFromString(dst)
	if err != nil {
		return err
	}

	dstHandle, err := windows.CreateFile(
		dstPath,
		windows.GENERIC_WRITE,
		windows.FILE_SHARE_WRITE,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	if err != nil {
		return err
	}
	defer windows.CloseHandle(dstHandle)

	// ファイルのタイムスタンプを設定
	err = windows.SetFileTime(dstHandle, &data.CreationTime, &data.LastAccessTime, &data.LastWriteTime)
	if err != nil {
		return err
	}

	return nil
}
