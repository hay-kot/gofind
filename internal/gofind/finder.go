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

	// Early termination for the current directory
	for _, dir := range dirs {
		for _, ignore := range fdr.Ignore {
			if Must(filepath.Match(ignore, dir.Name())) {
				return
			}
		}

		// Match the pattern
		if Must(filepath.Match(pattern, dir.Name())) {
			channel <- filepath.Join(root, dir.Name())
			return
		}
	}

	// Use a buffered channel as a semaphore to limit goroutines
	sem := make(chan struct{}, 10) // Max 10 concurrent goroutines

	var wg sync.WaitGroup
	for _, dir := range dirs {
		if dir.IsDir() {
			wg.Add(1)
			sem <- struct{}{} // Acquire semaphore

			go func(d fs.DirEntry) {
				defer wg.Done()
				defer func() { <-sem }() // Release semaphore
				fdr.DirWalker(channel, filepath.Join(root, d.Name()), pattern)
			}(dir)
		}
	}

	wg.Wait()
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

	// TODO: Do I need this or can I use the original wait group?
	collectWg.Wait()
	return matches, nil
}
