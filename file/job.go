package file

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type jobStatus struct {
	mu         sync.Mutex
	wg         sync.WaitGroup
	ch         chan string
	successCnt int32
	errorCnt   int32
	totalCnt   int32

	successFileCnt int32
	errorFileCnt   int32
	skipFileCnt    int32

	errorFiles []string
	errorDirs  []string	
}

func (j *jobStatus) getStatus() string {
	successCnt := atomic.LoadInt32(&j.successCnt)
	errCnt := atomic.LoadInt32(&j.errorCnt)
	progress := float32(successCnt+errCnt) / float32(j.totalCnt) * 100
	return fmt.Sprintf("[%3d%%] %d/%d(%d)",
		int(progress), successCnt+errCnt, j.totalCnt, errCnt,
	)
}

func (j *jobStatus) addSuccess() {
	atomic.AddInt32(&j.successCnt, 1)
}

func (j *jobStatus) addError() {
	atomic.AddInt32(&j.errorCnt, 1)
}

func (j *jobStatus) addSuccessFile() {
	atomic.AddInt32(&j.successFileCnt, 1)
}

func (j *jobStatus) addSkipFile() {
	atomic.AddInt32(&j.skipFileCnt, 1)
}

func (j *jobStatus) addErrorFile(file string) {
	atomic.AddInt32(&j.errorFileCnt, 1)

	j.mu.Lock()
	defer j.mu.Unlock()
	if j.errorFiles == nil {
		j.errorFiles = []string{}
	}
	j.errorFiles = append(j.errorFiles, file)
}

func (j *jobStatus) addErrorDirs(file string) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.errorDirs == nil {
		j.errorDirs = []string{}
	}
	j.errorFiles = append(j.errorDirs, file)
}
