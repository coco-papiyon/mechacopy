package directory

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/coco-papiyon/mechacopy/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDir = "testdata"

var recursionCase = []testutil.TestCase{
	{
		Name:     "empty",
		TestDirs: []string{"empty"},
	},
	{
		Name:     "single",
		TestDirs: []string{"single", "single/a"},
	},
	{
		Name:     "sub",
		TestDirs: []string{"sub", "sub/a", "sub/b", "sub/a/c"},
	},
}

func checkDirs(t *testing.T, tt testutil.TestCase, dirs []string, err error) {
	require.NoError(t, err, "getDirs")
	assert.Len(t, dirs, len(tt.TestDirs),
		"Dirs Count should be "+strconv.Itoa(len(tt.TestDirs)))
	for _, name := range tt.TestDirs {
		assert.Contains(t, dirs, name, "ディレクトリ存在チェック："+name)
	}
	for _, name := range tt.TestFiles {
		assert.NotContains(t, dirs, name, "ファイル存在チェック："+name)
	}
}

// ディレクトリ一覧取得
func TestGetSubDirs(t *testing.T) {
	tests := []testutil.TestCase{
		{
			Name: "empty",
		},
		{
			Name:     "single",
			TestDirs: []string{"a"},
		},
		{
			Name:      "sub",
			TestDirs:  []string{"a", "b"},
			TestFiles: []string{"a/file1.txt", "b/file2.txt"},
			ExtraDirs: []string{"a/c"},
		},
	}

	os.RemoveAll(testDir)

	// 指定ディレクトリが存在しない場合
	t.Run("Not Exist Directory", func(t *testing.T) {
		_, err := getSubDirs(testDir)
		assert.Error(t, err, "getDirs", testDir)
	})

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			defer os.RemoveAll(testDir)

			baseDir := filepath.Join(testDir, tt.Name)
			testutil.PrepareDirs(t, tt, baseDir)
			dirs, err := getSubDirs(baseDir)
			checkDirs(t, tt, dirs, err)
		})
	}
}

func TestGetDirRcrs(t *testing.T) {
	tests := recursionCase

	os.RemoveAll(testDir)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			defer os.RemoveAll(testDir)

			testutil.PrepareDirs(t, tt, testDir)
			dirs, err := getDirRecursion(testDir, tt.Name)
			checkDirs(t, tt, dirs, err)
		})
	}
}

func TestGetDir(t *testing.T) {
	tests := recursionCase

	os.RemoveAll(testDir)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			defer os.RemoveAll(testDir)

			testutil.PrepareDirs(t, tt, testDir)
			dirs, err := GetDirs(testDir)
			tt.TestDirs = append(tt.TestDirs, ".")
			checkDirs(t, tt, dirs, err)
		})
	}
}
