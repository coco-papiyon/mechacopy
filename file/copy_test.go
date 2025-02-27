package file

import (
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Infoログを非表示にする
func disableInfoLog() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})
	slog.SetDefault(slog.New(handler))
}

// Infoログを表示する
func enableInfoLog() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(handler))
}

func checkCopy(t *testing.T, src, dst string, err error) {
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

func TestCopyFile(t *testing.T) {
	tests := []testCase{
		{
			"dst_exist",
			[]string{"src", "dst"},
			[]string{"src/file1.txt"},
			[]string{},
			[]string{},
			nil,
		},
		{
			"dst_not_exist",
			[]string{"src"},
			[]string{"src/file1.txt"},
			[]string{},
			[]string{},
			nil,
		},
	}

	os.RemoveAll(testDir)

	// 指定ディレクトリが存在しない場合
	t.Run("Not Exist Directory", func(t *testing.T) {
		err := copyFile(testDir, "dst")
		assert.Error(t, err, "copyFile", testDir)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.RemoveAll(testDir)

			// 準備
			baseDir := filepath.Join(testDir, tt.name)
			prepareDirs(t, tt, baseDir)
			src := filepath.Join(baseDir, tt.testFiles[0])
			dst := filepath.Base(tt.testFiles[0])
			dst = filepath.Join(baseDir, "dst", dst)

			// コピー実行
			err := copyFile(src, dst)
			checkCopy(t, src, dst, err)
		})
	}
}

func TestCopyFiles(t *testing.T) {
	tests := []testCase{
		{
			"files",
			[]string{},
			[]string{"file1.txt", "file2.txt"},
			[]string{},
			[]string{},
			nil,
		},
		{
			"filesWithPtn",
			[]string{},
			[]string{"file1.txt"},
			[]string{},
			[]string{"file2.log"},
			[]string{"*.txt"},
		},
	}

	os.RemoveAll(testDir)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.RemoveAll(testDir)

			// 準備
			baseDir := filepath.Join(testDir, tt.name)
			srcDir := filepath.Join(baseDir, "src")
			prepareDirs(t, tt, srcDir)

			// コピー実行
			job := &jobStatus{}
			err := copyFiles(baseDir, "src", "testdata/dst", job)
			for _, file := range tt.testFiles {
				src := filepath.Join(srcDir, file)
				dst := filepath.Base(file)
				dst = filepath.Join(testDir, "dst", dst)
				checkCopy(t, src, dst, err)
			}
		})
	}
}

func TestCopyFilesParallel(t *testing.T) {
	tests := []testCase{
		{
			"files",
			[]string{"1", "2", "1/3"},
			[]string{"file1.txt", "file2.txt",
				"1/file1.txt", "1/file2.txt",
				"2/file1.txt", "2/file2.txt",
				"1/3/file1.txt", "1/3/file2.txt",
			},
			[]string{},
			[]string{},
			nil,
		},
	}

	testMany := testCase{}
	for i := range 10 {
		vali := strconv.Itoa(i)
		for j := range 10 {
			valj := strconv.Itoa(j)
			testMany.testDirs = append(testMany.testDirs, vali+"/"+valj)
			testMany.testFiles = append(testMany.testFiles, vali+"/"+valj+".txt")
			for k := range 10 {
				valk := strconv.Itoa(k)
				testMany.testFiles = append(testMany.testFiles, vali+"/"+valj+"/"+valk+".txt")
			}
		}
	}
	tests = append(tests, testMany)

	os.RemoveAll(testDir)

	// Infoログを非表示にする(テスト後に戻す)
	disableInfoLog()
	defer enableInfoLog()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.RemoveAll(testDir)

			// 準備
			baseDir := filepath.Join(testDir, tt.name, "src")
			prepareDirs(t, tt, baseDir)
			// コピー実行
			err := CopyFiles(baseDir, "testdata/dst")
			for _, file := range tt.testFiles {
				src := filepath.Join(baseDir, file)
				dst := filepath.Join(testDir, "dst", file)
				checkCopy(t, src, dst, err)
			}
		})
	}
}

func TestIsCopyFile(t *testing.T) {
	tests := []struct {
		target  []string
		name    string
		matched bool
	}{
		{nil, "file1", true},
		{nil, "file1.txt", true},
		{nil, "test/file1", true},
		{[]string{"*.txt"}, "file1.txt", true},
		{[]string{"*.*"}, "file1.txt", true},
		{[]string{"*.*"}, "file1", false},
		{[]string{"*.txt"}, "file1", false},
		{[]string{"file*"}, "file1", true},
		{[]string{"file*"}, "afile", false},
		{[]string{"*.log", "*.txt"}, "file1.txt", true},
		{[]string{"*.log", "*.txt"}, "file1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.target != nil {
				Config.TargetFiles = tt.target
			}
			matched := isCopyFile(tt.name)
			assert.Equal(t, tt.matched, matched, "isCopyFile %v, %s", tt.target, tt.name)
		})
	}
}

func TestIsFileDiff(t *testing.T) {
	os.RemoveAll(testDir)
	defer os.RemoveAll(testDir)

	err := os.MkdirAll(testDir, 0755)
	require.NoError(t, err, "準備：ディレクトリ作成 %s", testDir)

	file1 := filepath.Join(testDir, "file1")
	file2 := filepath.Join(testDir, "file2_____")
	file3 := filepath.Join(testDir, "file3__________")
	file4 := filepath.Join(testDir, "file4")

	err = createTestFile(file1)
	require.NoError(t, err, "準備：ファイル作成 %s", file1)
	time.Sleep(100 * time.Millisecond)
	err = createTestFile(file2)
	require.NoError(t, err, "準備：ファイル作成 %s", file2)
	time.Sleep(100 * time.Millisecond)
	copyFile(file1, file3)
	require.NoError(t, err, "準備：コピーファイル %s %s", file1, file3)
	time.Sleep(100 * time.Millisecond)
	err = createTestFile(file4)
	require.NoError(t, err, "準備：ファイル作成 %s", file4)

	diff := isFileDiff(file1, file2)
	assert.Equal(t, true, diff, "different file %s %s", file1, file2)

	diff = isFileDiff(file1, "aaa")
	assert.Equal(t, true, diff, "dst file is not exist %s %s", file1, "aaa")

	diff = isFileDiff(file1, file3)
	assert.Equal(t, false, diff, "same and same file %s %s", file1, file3)

	diff = isFileDiff(file1, file4)
	assert.Equal(t, false, diff, "same and new file %s %s", file1, file4)

	diff = isFileDiff(file4, file1)
	assert.Equal(t, true, diff, "same and old file %s %s", file4, file1)

	// fileInfo1, _ := os.Stat(file1)
	// fileInfo2, _ := os.Stat(file2)
	// fileInfo3, _ := os.Stat(file3)
	// fileInfo4, _ := os.Stat(file4)
	// fmt.Println(fileInfo1.ModTime().UnixNano(), fileInfo1.Name(), fileInfo1.Size())
	// fmt.Println(fileInfo2.ModTime().UnixNano(), fileInfo2.Name(), fileInfo2.Size())
	// fmt.Println(fileInfo3.ModTime().UnixNano(), fileInfo3.Name(), fileInfo3.Size())
	// fmt.Println(fileInfo4.ModTime().UnixNano(), fileInfo4.Name(), fileInfo4.Size())
}

func TestCopyRetry(t *testing.T) {
	os.RemoveAll(testDir)
	defer os.RemoveAll(testDir)

	// 準備
	srcDir := filepath.Join(testDir, "src")
	dstDir := filepath.Join(testDir, "dst")
	errFiles := []string{"file1.txt"}
	job := &jobStatus{
		errorFiles: errFiles,
	}
	// 数秒後にファイルを作成
	go func() {
		time.Sleep(3 * time.Second)
		for _, file := range errFiles {
			testfile := filepath.Join(srcDir, file)
			createTestFile(testfile)
		}
	}()

	// コピー実行
	Config.RetryCount = 10
	Config.SleepTime = 1
	copyFileTry(srcDir, dstDir, job)

	src := filepath.Join(srcDir, "file1.txt")
	dst := filepath.Join(dstDir, "file1.txt")
	checkCopy(t, src, dst, nil)
}
