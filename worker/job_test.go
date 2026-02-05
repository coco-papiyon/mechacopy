package worker

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJobStatus_GetStatus(t *testing.T) {
	tests := []struct {
		name       string
		successCnt int32
		errorCnt   int32
		totalCnt   int32
		want       string
	}{
		{
			name:       "0%進捗",
			successCnt: 0,
			errorCnt:   0,
			totalCnt:   100,
			want:       "[  0%] 0/100(0)",
		},
		{
			name:       "50%進捗",
			successCnt: 50,
			errorCnt:   0,
			totalCnt:   100,
			want:       "[ 50%] 50/100(0)",
		},
		{
			name:       "100%進捗",
			successCnt: 100,
			errorCnt:   0,
			totalCnt:   100,
			want:       "[100%] 100/100(0)",
		},
		{
			name:       "エラーあり",
			successCnt: 90,
			errorCnt:   10,
			totalCnt:   100,
			want:       "[100%] 100/100(10)",
		},
		{
			name:       "小数点切り捨て",
			successCnt: 1,
			errorCnt:   0,
			totalCnt:   3,
			want:       "[ 33%] 1/3(0)",
		},
		{
			name:       "totalが1",
			successCnt: 1,
			errorCnt:   0,
			totalCnt:   1,
			want:       "[100%] 1/1(0)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JobStatus{
				successCnt: tt.successCnt,
				errorCnt:   tt.errorCnt,
				totalCnt:   tt.totalCnt,
			}
			got := j.GetStatus()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestJobStatus_AddSuccess(t *testing.T) {
	tests := []struct {
		name     string
		initial  int32
		addCount int
		want     int32
	}{
		{
			name:     "1回加算",
			initial:  0,
			addCount: 1,
			want:     1,
		},
		{
			name:     "複数回加算",
			initial:  5,
			addCount: 10,
			want:     15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JobStatus{
				successCnt: tt.initial,
			}
			for i := 0; i < tt.addCount; i++ {
				j.AddSuccess()
			}
			assert.Equal(t, tt.want, j.successCnt)
		})
	}
}

func TestJobStatus_AddSuccess_Concurrent(t *testing.T) {
	job := &JobStatus{}
	var wg sync.WaitGroup
	goroutines := 100
	iterations := 100

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				job.AddSuccess()
			}
		}()
	}
	wg.Wait()

	expected := int32(goroutines * iterations)
	assert.Equal(t, expected, job.successCnt)
}

func TestJobStatus_AddError(t *testing.T) {
	tests := []struct {
		name     string
		initial  int32
		addCount int
		want     int32
	}{
		{
			name:     "1回加算",
			initial:  0,
			addCount: 1,
			want:     1,
		},
		{
			name:     "複数回加算",
			initial:  3,
			addCount: 7,
			want:     10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JobStatus{
				errorCnt: tt.initial,
			}
			for i := 0; i < tt.addCount; i++ {
				j.AddError()
			}
			assert.Equal(t, tt.want, j.errorCnt)
		})
	}
}

func TestJobStatus_AddError_Concurrent(t *testing.T) {
	job := &JobStatus{}
	var wg sync.WaitGroup
	goroutines := 100
	iterations := 100

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				job.AddError()
			}
		}()
	}
	wg.Wait()

	expected := int32(goroutines * iterations)
	assert.Equal(t, expected, job.errorCnt)
}

func TestJobStatus_AddSuccessFile(t *testing.T) {
	tests := []struct {
		name     string
		initial  int32
		addCount int
		want     int32
	}{
		{
			name:     "1回加算",
			initial:  0,
			addCount: 1,
			want:     1,
		},
		{
			name:     "複数回加算",
			initial:  10,
			addCount: 20,
			want:     30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JobStatus{
				successFileCnt: tt.initial,
			}
			for i := 0; i < tt.addCount; i++ {
				j.AddSuccessFile()
			}
			assert.Equal(t, tt.want, j.successFileCnt)
		})
	}
}

func TestJobStatus_AddSuccessFile_Concurrent(t *testing.T) {
	job := &JobStatus{}
	var wg sync.WaitGroup
	goroutines := 100
	iterations := 100

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				job.AddSuccessFile()
			}
		}()
	}
	wg.Wait()

	expected := int32(goroutines * iterations)
	assert.Equal(t, expected, job.successFileCnt)
}

func TestJobStatus_AddSkipFile(t *testing.T) {
	tests := []struct {
		name     string
		initial  int32
		addCount int
		want     int32
	}{
		{
			name:     "1回加算",
			initial:  0,
			addCount: 1,
			want:     1,
		},
		{
			name:     "複数回加算",
			initial:  5,
			addCount: 5,
			want:     10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JobStatus{
				skipFileCnt: tt.initial,
			}
			for i := 0; i < tt.addCount; i++ {
				j.AddSkipFile()
			}
			assert.Equal(t, tt.want, j.skipFileCnt)
		})
	}
}

func TestJobStatus_AddSkipFile_Concurrent(t *testing.T) {
	job := &JobStatus{}
	var wg sync.WaitGroup
	goroutines := 100
	iterations := 100

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				job.AddSkipFile()
			}
		}()
	}
	wg.Wait()

	expected := int32(goroutines * iterations)
	assert.Equal(t, expected, job.skipFileCnt)
}

func TestJobStatus_AddErrorFile(t *testing.T) {
	tests := []struct {
		name       string
		files      []string
		wantLen    int
		wantFiles  []string
	}{
		{
			name:      "1ファイル追加",
			files:     []string{"file1.txt"},
			wantLen:   1,
			wantFiles: []string{"file1.txt"},
		},
		{
			name:      "複数ファイル追加",
			files:     []string{"file1.txt", "file2.txt", "file3.txt"},
			wantLen:   3,
			wantFiles: []string{"file1.txt", "file2.txt", "file3.txt"},
		},
		{
			name:      "空文字列追加",
			files:     []string{""},
			wantLen:   1,
			wantFiles: []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JobStatus{}
			for _, file := range tt.files {
				j.AddErrorFile(file)
			}
			assert.Equal(t, tt.wantLen, len(j.errorFiles))
			assert.Equal(t, tt.wantFiles, j.errorFiles)
		})
	}
}

func TestJobStatus_AddErrorFile_Concurrent(t *testing.T) {
	j := &JobStatus{}
	var wg sync.WaitGroup
	goroutines := 100

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		i := i
		go func() {
			defer wg.Done()
			j.AddErrorFile("file_" + string(rune(i)))
		}()
	}
	wg.Wait()

	assert.Equal(t, goroutines, len(j.errorFiles))
}

func TestJobStatus_AddErrorDirs(t *testing.T) {
	tests := []struct {
		name      string
		dirs      []string
		wantLen   int
	}{
		{
			name:    "1ディレクトリ追加",
			dirs:    []string{"dir1"},
			wantLen: 0,
		},
		{
			name:    "複数ディレクトリ追加",
			dirs:    []string{"dir1", "dir2", "dir3"},
			wantLen: 0,
		},
		{
			name:    "空文字列追加",
			dirs:    []string{""},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &JobStatus{}
			for _, dir := range tt.dirs {
				j.AddErrorDirs(dir)
			}
			assert.Equal(t, tt.wantLen, len(j.errorDirs))
		})
	}
}

func TestJobStatus_AddErrorDirs_Concurrent(t *testing.T) {
	job := &JobStatus{}
	var wg sync.WaitGroup
	goroutines := 100

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		i := i
		go func() {
			defer wg.Done()
			job.AddErrorDirs("dir_" + string(rune(i)))
		}()
	}
	wg.Wait()

	assert.Equal(t, 0, len(job.errorDirs))
}
