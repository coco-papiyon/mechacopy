package worker

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/coco-papiyon/mechacopy/filecopy"
	"github.com/coco-papiyon/mechacopy/testutil"
	"github.com/stretchr/testify/assert"
)

const testDir = "testdata"

var testConfig = InitConfig()

func TestCopyFilesParallel(t *testing.T) {
	tests := []testutil.TestCase{
		{
			Name:     "files",
			TestDirs: []string{"1", "2", "1/3"},
			TestFiles: []string{"file1.txt", "file2.txt",
				"1/file1.txt", "1/file2.txt",
				"2/file1.txt", "2/file2.txt",
				"1/3/file1.txt", "1/3/file2.txt",
			},
			ExtraDirs:   []string{},
			ExtraFiles:  []string{},
			TargetFiles: nil,
		},
	}

	testMany := testutil.TestCase{}
	for i := range 10 {
		vali := strconv.Itoa(i)
		for j := range 10 {
			valj := strconv.Itoa(j)
			testMany.TestDirs = append(testMany.TestDirs, vali+"/"+valj)
			testMany.TestFiles = append(testMany.TestFiles, vali+"/"+valj+".txt")
			for k := range 10 {
				valk := strconv.Itoa(k)
				testMany.TestFiles = append(testMany.TestFiles, vali+"/"+valj+"/"+valk+".txt")
			}
		}
	}
	tests = append(tests, testMany)

	os.RemoveAll(testDir)

	// Infoログを非表示にする(テスト後に戻す)
	testutil.DisableInfoLog()
	defer testutil.EnableInfoLog()

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			defer os.RemoveAll(testDir)

			// 準備
			baseDir := filepath.Join(testDir, tt.Name, "src")
			testutil.PrepareDirs(t, tt, baseDir)

			var runner Runner = &CopyRunner{
				Destination: "testdata/dst",
			}

			// コピー実行
			err := RunMecha(baseDir, runner, testConfig)
			for _, file := range tt.TestFiles {
				src := filepath.Join(baseDir, file)
				dst := filepath.Join(testDir, "dst", file)
				testutil.CheckCopy(t, src, dst, err)
			}
		})
	}
}

func TestCopyFiles(t *testing.T) {
	tests := []testutil.TestCase{
		{
			Name:        "files",
			TestDirs:    []string{},
			TestFiles:   []string{"file1.txt", "file2.txt"},
			ExtraDirs:   []string{},
			ExtraFiles:  []string{},
			TargetFiles: nil,
		},
		{
			Name:        "filesWithPtn",
			TestDirs:    []string{},
			TestFiles:   []string{"file1.txt"},
			ExtraDirs:   []string{},
			ExtraFiles:  []string{"file2.log"},
			TargetFiles: []string{"*.txt"},
		},
	}

	os.RemoveAll(testDir)

	// Infoログを非表示にする(テスト後に戻す)
	testutil.DisableInfoLog()
	defer testutil.EnableInfoLog()

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			defer os.RemoveAll(testDir)

			// 準備
			baseDir := filepath.Join(testDir, tt.Name)
			srcDir := filepath.Join(baseDir, "src")
			dstDir := filepath.Join(baseDir, "dst")
			testutil.PrepareDirs(t, tt, srcDir)

			var runner Runner = &CopyRunner{
				Destination: dstDir,
			}

			// コピー実行
			err := RunMecha(srcDir, runner, testConfig)
			for _, file := range tt.TestFiles {
				src := filepath.Join(srcDir, file)
				dst := filepath.Join(dstDir, file)
				testutil.CheckCopy(t, src, dst, err)
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
				testConfig.TargetFiles = tt.target
			}
			matched := filecopy.IsCopyFile(tt.name, testConfig.TargetFiles)
			assert.Equal(t, tt.matched, matched, "isCopyFile %v, %s", tt.target, tt.name)
		})
	}
}

func TestCopyRetry(t *testing.T) {
	os.RemoveAll(testDir)
	defer os.RemoveAll(testDir)

	// 準備
	srcDir := filepath.Join(testDir, "src")
	dstDir := filepath.Join(testDir, "dst")
	errFiles := []string{}
	for i := 0; i < 30; i++ {
		errFiles = append(errFiles, fmt.Sprintf("test%d.txt", i))
	}

	job := &JobStatus{
		errorFiles: errFiles,
		config:     testConfig,
	}
	// 数秒後にファイルを作成
	go func() {
		time.Sleep(2 * time.Second)
		for i, file := range errFiles {
			time.Sleep(time.Duration(i*10) * time.Millisecond)
			testfile := filepath.Join(srcDir, file)
			testutil.CreateTestFile(testfile)
		}
	}()

	// コピー実行
	var runner Runner = &CopyRunner{
		Destination: dstDir,
	}

	testConfig.RetryCount = 10
	testConfig.SleepTime = 1
	runRetry(srcDir, runner, job)

	for _, file := range errFiles {
		src := filepath.Join(srcDir, file)
		dst := filepath.Join(dstDir, file)
		testutil.CheckCopy(t, src, dst, nil)
	}
}
