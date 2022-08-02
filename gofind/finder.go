package gofind

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

type Finder struct {
	MaxRecursion     int
	CurrentRecursion int
	Ignore           []string
}

type Match struct {
	Name string
	Path string
}

func (fdr *Finder) DirWalker(channel chan string, root string, pattern string) {
	dirs := Must(os.ReadDir(root))
	w := sync.WaitGroup{}

	found := false
	collection := []fs.DirEntry{}

	for _, dir := range dirs {
		// Early Return if directory is in ignore list
		for _, ignore := range fdr.Ignore {
			if Must(filepath.Match(ignore, dir.Name())) {
				return
			}
		}

		if Must(filepath.Match(pattern, dir.Name())) {
			channel <- filepath.Join(root, dir.Name())
			found = true
		} else if dir.IsDir() {
			func(dir fs.DirEntry) {
				collection = append(collection, dir)
			}(dir)
		}
	}

	if !found {
		for _, dir := range collection {
			w.Add(1)
			go func(d fs.DirEntry) {
				defer w.Done()
				fdr.DirWalker(channel, filepath.Join(root, d.Name()), pattern)
			}(dir)
		}
	}

	w.Wait()
}

func (fdr *Finder) FindAll(wg *sync.WaitGroup, channel chan string, root string, pattern string) {
	defer wg.Done()

	fdr.DirWalker(channel, root, pattern)

}

func (fdr *Finder) CollectResults(wg *sync.WaitGroup, channel chan string, results *[]string) {
	defer wg.Done()

	for result := range channel {
		*results = append(*results, result)
	}

}

func (fdr *Finder) Find(path []string, glob string) ([]string, error) {
	var (
		matches = []string{}
		results = make(chan string)
	)

	finderWg := sync.WaitGroup{}
	finderWg.Add(len(path))

	for _, root := range path {
		go fdr.FindAll(&finderWg, results, root, glob)
	}

	collectWg := sync.WaitGroup{}
	collectWg.Add(1)

	go fdr.CollectResults(&collectWg, results, &matches)

	finderWg.Wait()
	close(results)

	collectWg.Wait()
	return matches, nil
}
