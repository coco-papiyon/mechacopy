package worker

import (
	"log/slog"
	"os"
	"path/filepath"
)

type DeleteRunner struct {
}

func (r DeleteRunner) Run(baseDir, srcDir string, job *JobStatus) error {
	target := filepath.Join(baseDir, srcDir)
	err := os.RemoveAll(target)
	if err != nil {
		slog.Error("Delete Directory", "directory", target, "ERROR", err)
		//job.AddErrorFile(target)
	}
	return nil
}

func (r DeleteRunner) Retry(srcDir, targetFile string) error {
	return nil
}
