//go:build !windows

package file

import "fmt"

func copyTimestampsWindows(src, dst string) error {
	return fmt.Errorf("OS Not Supported")
}
