package testutil

import (
	"log/slog"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	Name        string
	TestDirs    []string
	TestFiles   []string
	ExtraDirs   []string
	ExtraFiles  []string
	TargetFiles []string
}

var randSource = rand.New(rand.NewSource(time.Now().UnixNano()))

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randString() string {
	n := randSource.Intn(100) + 1
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// テスト用のファイル作成
func CreateTestFile(filePath string) error {

	dirpath := filepath.Dir(filePath)
	err := os.MkdirAll(dirpath, 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	randomString := randString()
	_, err = file.WriteString(filePath + randomString)
	if err != nil {
		return err
	}
	return nil
}

func PrepareDirs(t *testing.T, tt TestCase, baseDir string) {
	err := os.MkdirAll(baseDir, 0755)
	require.NoError(t, err, "準備：ディレクトリ作成")

	// 準備
	for _, name := range tt.TestDirs {
		path := filepath.Join(baseDir, name)
		err := os.MkdirAll(path, 0755)
		require.NoError(t, err, "準備：ディレクトリ作成")
	}
	for _, name := range tt.ExtraDirs {
		path := filepath.Join(baseDir, name)
		err := os.MkdirAll(path, 0755)
		require.NoError(t, err, "準備：ディレクトリ作成")
	}
	for _, name := range tt.TestFiles {
		path := filepath.Join(baseDir, name)
		err := CreateTestFile(path)
		require.NoError(t, err, "準備：ファイル作成")
	}
	for _, name := range tt.ExtraFiles {
		path := filepath.Join(baseDir, name)
		err := CreateTestFile(path)
		require.NoError(t, err, "準備：ファイル作成")
	}
}

// Infoログを非表示にする
func setLogLevel(level slog.Level) {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	slog.SetDefault(slog.New(handler))
}

// Infoログを非表示にする
func DisableInfoLog() {
	setLogLevel(slog.LevelWarn)
}

// Infoログを表示する
func EnableInfoLog() {
	setLogLevel(slog.LevelInfo)
}

func CheckCopy(t *testing.T, src, dst string, err error) {
	require.NoError(t, err, "copyFile")

	srcInfo, err := os.Stat(src)
	require.NoError(t, err, "Stat %s", src)

	srcData, err := os.ReadFile(src)
	require.NoError(t, err, "ReadFile %s", src)

	dstInfo, err := os.Stat(dst)
	require.NoError(t, err, "Stat %s", dst)

	dstData, err := os.ReadFile(dst)
	require.NoError(t, err, "ReadFile %s", dst)

	assert.Equal(t, srcData, dstData, "Copy Check")
	assert.Equal(t, srcInfo.ModTime().UnixMilli(), dstInfo.ModTime().UnixMilli(), "Copy Check ModTime")
}
