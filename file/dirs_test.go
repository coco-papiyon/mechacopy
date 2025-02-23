package file

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDir = "testdata"

type testCase struct {
	name        string
	testDirs    []string
	testFiles   []string
	extraDirs   []string
	extraFiles  []string
	targetFiles []string
}

var recursionCase = []testCase{
	{
		"empty",
		[]string{"empty"},
		[]string{},
		[]string{},
		[]string{},
		nil,
	},
	{
		"single",
		[]string{"single", "single/a"},
		[]string{},
		[]string{},
		[]string{},
		nil,
	},
	{
		"sub",
		[]string{"sub", "sub/a", "sub/b", "sub/a/c"},
		[]string{},
		[]string{},
		[]string{},
		nil,
	},
}

// テスト用のファイル作成
func createTestFile(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(filePath)
	if err != nil {
		return err
	}
	return nil
}

func prepareDirs(t *testing.T, tt testCase, baseDir string) {
	err := os.MkdirAll(baseDir, 0755)
	require.NoError(t, err, "準備：ディレクトリ作成")

	// 準備
	for _, name := range tt.testDirs {
		path := filepath.Join(baseDir, name)
		err := os.MkdirAll(path, 0755)
		require.NoError(t, err, "準備：ディレクトリ作成")
	}
	for _, name := range tt.extraDirs {
		path := filepath.Join(baseDir, name)
		err := os.MkdirAll(path, 0755)
		require.NoError(t, err, "準備：ディレクトリ作成")
	}
	for _, name := range tt.testFiles {
		path := filepath.Join(baseDir, name)
		err := createTestFile(path)
		require.NoError(t, err, "準備：ファイル作成")
	}
	for _, name := range tt.extraFiles {
		path := filepath.Join(baseDir, name)
		err := createTestFile(path)
		require.NoError(t, err, "準備：ファイル作成")
	}
}

func checkDirs(t *testing.T, tt testCase, dirs []string, err error) {
	require.NoError(t, err, "getDirs")
	assert.Len(t, dirs, len(tt.testDirs),
		"Dirs Count should be "+strconv.Itoa(len(tt.testDirs)))
	for _, name := range tt.testDirs {
		assert.Contains(t, dirs, name, "ディレクトリ存在チェック："+name)
	}
	for _, name := range tt.testFiles {
		assert.NotContains(t, dirs, name, "ファイル存在チェック："+name)
	}
}

// ディレクトリ一覧取得
func TestGetSubDirs(t *testing.T) {
	tests := []testCase{
		{
			"empty",
			[]string{},
			[]string{},
			[]string{},
			[]string{},
			nil,
		},
		{
			"single",
			[]string{"a"},
			[]string{},
			[]string{},
			[]string{},
			nil,
		},
		{
			"sub",
			[]string{"a", "b"},
			[]string{"a/file1.txt", "b/file2.txt"},
			[]string{"a/c"},
			[]string{},
			nil,
		},
	}

	os.RemoveAll(testDir)

	// 指定ディレクトリが存在しない場合
	t.Run("Not Exist Directory", func(t *testing.T) {
		_, err := getSubDirs(testDir)
		assert.Error(t, err, "getDirs", testDir)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.RemoveAll(testDir)

			baseDir := filepath.Join(testDir, tt.name)
			prepareDirs(t, tt, baseDir)
			dirs, err := getSubDirs(baseDir)
			checkDirs(t, tt, dirs, err)
		})
	}
}

func TestGetDirRcrs(t *testing.T) {
	tests := recursionCase

	os.RemoveAll(testDir)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.RemoveAll(testDir)

			prepareDirs(t, tt, testDir)
			dirs, err := getDirRecursion(testDir, tt.name)
			checkDirs(t, tt, dirs, err)
		})
	}
}

func TestGetDir(t *testing.T) {
	tests := recursionCase

	os.RemoveAll(testDir)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.RemoveAll(testDir)

			prepareDirs(t, tt, testDir)
			dirs, err := GetDirs(testDir)
			tt.testDirs = append(tt.testDirs, ".")
			checkDirs(t, tt, dirs, err)
		})
	}
}
