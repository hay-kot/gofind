package gofind

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

func DirWalker(channel chan string, root string, pattern string) {
	dirs := Must(os.ReadDir(root))

	w := sync.WaitGroup{}

	children := make(chan fs.DirEntry)
	found := false

	collection := []fs.DirEntry{}

	for _, dir := range dirs {
		func(dir fs.DirEntry) {
			if Must(filepath.Match(pattern, dir.Name())) {
				channel <- filepath.Join(root, dir.Name())
				found = true
			}
			if dir.IsDir() {
				func(dir fs.DirEntry) {
					collection = append(collection, dir)
				}(dir)
			}
		}(dir)
	}

	w.Wait()
	close(children)

	if !found {
		for _, dir := range collection {
			w.Add(1)
			go func(d fs.DirEntry) {
				defer w.Done()
				DirWalker(channel, filepath.Join(root, d.Name()), pattern)
			}(dir)
		}
	}

	w.Wait()
}

func FindAll(wg *sync.WaitGroup, channel chan string, root string, pattern string) {
	defer wg.Done()
	defer close(channel)

	DirWalker(channel, root, pattern)

}

func CollectResults(wg *sync.WaitGroup, channel chan string, results *[]string) {
	defer wg.Done()

	for result := range channel {
		*results = append(*results, result)
	}

}

func Finder(path string, glob string) ([]string, error) {
	matches := []string{}

	wg := sync.WaitGroup{}

	results := make(chan string)

	wg.Add(1)
	go FindAll(&wg, results, path, glob)
	wg.Add(1)
	go CollectResults(&wg, results, &matches)
	wg.Wait()

	return matches, nil
}
