package filecopy

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/coco-papiyon/mechacopy/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDir = "testdata"

func TestCopyFile(t *testing.T) {
	tests := []testutil.TestCase{
		{
			Name:        "dst_exist",
			TestDirs:    []string{"src", "dst"},
			TestFiles:   []string{"src/file1.txt"},
			ExtraDirs:   []string{},
			ExtraFiles:  []string{},
			TargetFiles: nil,
		},
		{
			Name:        "dst_not_exist",
			TestDirs:    []string{"src"},
			TestFiles:   []string{"src/file1.txt"},
			ExtraDirs:   []string{},
			ExtraFiles:  []string{},
			TargetFiles: nil,
		},
	}

	os.RemoveAll(testDir)

	// 指定ディレクトリが存在しない場合
	t.Run("Not Exist Directory", func(t *testing.T) {
		err := CopyFile(testDir, "dst")
		assert.Error(t, err, "copyFile", testDir)
	})

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			defer os.RemoveAll(testDir)

			// 準備
			baseDir := filepath.Join(testDir, tt.Name)
			testutil.PrepareDirs(t, tt, baseDir)
			src := filepath.Join(baseDir, tt.TestFiles[0])
			dst := filepath.Base(tt.TestFiles[0])
			dst = filepath.Join(baseDir, "dst", dst)

			// コピー実行
			err := CopyFile(src, dst)
			testutil.CheckCopy(t, src, dst, err)
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

	err = testutil.CreateTestFile(file1)
	require.NoError(t, err, "準備：ファイル作成 %s", file1)
	time.Sleep(100 * time.Millisecond)
	err = testutil.CreateTestFile(file2)
	require.NoError(t, err, "準備：ファイル作成 %s", file2)
	time.Sleep(100 * time.Millisecond)
	CopyFile(file1, file3)
	require.NoError(t, err, "準備：コピーファイル %s %s", file1, file3)
	time.Sleep(100 * time.Millisecond)
	err = copyData(file1, file4)
	require.NoError(t, err, "準備：ファイル作成 %s", file4)

	diff := IsFileDiff(file1, file2)
	assert.Equal(t, true, diff, "different file %s %s", file1, file2)

	diff = IsFileDiff(file1, "aaa")
	assert.Equal(t, true, diff, "dst file is not exist %s %s", file1, "aaa")

	diff = IsFileDiff(file1, file3)
	assert.Equal(t, false, diff, "same and same file %s %s", file1, file3)

	diff = IsFileDiff(file1, file4)
	assert.Equal(t, false, diff, "same and new file %s %s", file1, file4)

	diff = IsFileDiff(file4, file1)
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
