package worker

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type JobStatus struct {
	mu sync.Mutex
	wg sync.WaitGroup
	ch chan string

	config *Config

	successCnt int32
	errorCnt   int32
	totalCnt   int32

	successFileCnt int32
	skipFileCnt    int32

	errorFiles []string
	errorDirs  []string
}

func (j *JobStatus) GetStatus() string {
	successCnt := atomic.LoadInt32(&j.successCnt)
	errCnt := atomic.LoadInt32(&j.errorCnt)
	progress := float32(successCnt+errCnt) / float32(j.totalCnt) * 100
	return fmt.Sprintf("[%3d%%] %d/%d(%d)",
		int(progress), successCnt+errCnt, j.totalCnt, errCnt,
	)
}

func (j *JobStatus) AddSuccess() {
	atomic.AddInt32(&j.successCnt, 1)
}

func (j *JobStatus) AddError() {
	atomic.AddInt32(&j.errorCnt, 1)
}

func (j *JobStatus) AddSuccessFile() {
	atomic.AddInt32(&j.successFileCnt, 1)
}

func (j *JobStatus) AddSkipFile() {
	atomic.AddInt32(&j.skipFileCnt, 1)
}

func (j *JobStatus) AddErrorFile(file string) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.errorFiles == nil {
		j.errorFiles = []string{}
	}
	j.errorFiles = append(j.errorFiles, file)
}

func (j *JobStatus) AddErrorDirs(file string) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.errorDirs == nil {
		j.errorDirs = []string{}
	}
	j.errorFiles = append(j.errorDirs, file)
}
