package directory

import (
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// 再帰的にディレクトリを取得(同時実行制御)
func getDirWorker(basedir, dirname string, wg *sync.WaitGroup, ch chan<- []string) {
	defer wg.Done()
	dir, err := getDirRecursion(basedir, dirname)
	if err != nil {
		slog.Error("ディレクトリ取得エラー", "ERROR", err, "Directory", dirname)
		return
	}
	ch <- dir
}

// 再帰的にディレクトリを取得
func getDirRecursion(basedir, dirname string) ([]string, error) {
	dirs := []string{dirname}

	// 子ディレクトリ一覧を取得
	dirpath := filepath.Join(basedir, dirname)
	childs, err := getSubDirs(dirpath)
	if err != nil {
		return dirs, err
	}

	// 子ディレクトリに対して再帰的にディレクトリ検索を行う
	for _, child := range childs {
		childDir := filepath.Join(dirname, child)
		grands, err := getDirRecursion(basedir, childDir)
		if err != nil {
			return dirs, err
		}
		dirs = append(dirs, grands...)
	}
	return dirs, nil
}

// ディレクトリ一覧を取得
func getSubDirs(path string) ([]string, error) {
	dirs := []string{}
	entries, err := os.ReadDir(path)
	if err != nil {
		return dirs, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}
	return dirs, nil
}

// ディレクトリ一覧を取得
func GetDirs(path string) ([]string, error) {
	dirs := []string{"."}

	// 直下のディレクトリを取得
	subDirs, err := getSubDirs(path)
	if err != nil {
		return dirs, err
	}

	// サブディレクトリごとに再帰的にディレクトリを取得
	var wg sync.WaitGroup
	ch := make(chan []string)
	for _, entry := range subDirs {
		wg.Add(1)
		go getDirWorker(path, entry, &wg, ch)
	}

	// 処理待ち
	go func() {
		wg.Wait()
		close(ch)
	}()

	// 結果を受け取る
	for result := range ch {
		dirs = append(dirs, result...)
	}

	// 降順にソート
	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i] > dirs[j]
	})

	return dirs, nil
}
