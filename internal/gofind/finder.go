package gofind

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/rs/zerolog/log"
)

const numWorkers = 32

type Finder struct {
	Ignore []string
}

type Match struct {
	Name string
	Path string
}

// isIgnored reports whether name matches any ignore pattern.
func (fdr *Finder) isIgnored(name string) bool {
	for _, pat := range fdr.Ignore {
		matched, err := filepath.Match(pat, name)
		if err != nil {
			log.Warn().Err(err).Str("pattern", pat).Str("name", name).Msg("invalid ignore pattern")
			continue
		}
		if matched {
			return true
		}
	}
	return false
}

// Find walks each root directory concurrently using a fixed worker pool and
// returns the paths of all entries whose name matches glob.
//
// A directory is not recursed into once a match is found within it.
func (fdr *Finder) Find(roots []string, glob string) ([]string, error) {
	// queue holds directories waiting to be scanned. The buffer is sized
	// generously to avoid blocking workers when they submit child directories.
	queue := make(chan string, 65536)
	results := make(chan string, 256)

	// wg tracks the number of directories that have been submitted but not yet
	// fully processed (including their children being submitted). wg.Wait()
	// signals that all traversal is complete.
	var wg sync.WaitGroup

	// submit enqueues a directory for processing. It increments wg before
	// pushing so wg.Wait() cannot return while the item is still pending.
	submit := func(dir string) {
		wg.Add(1)
		select {
		case queue <- dir:
		default:
			// Buffer full — push in a goroutine to avoid blocking the worker.
			// wg is already incremented, so this item is tracked.
			go func() { queue <- dir }()
		}
	}

	// Seed with the initial roots.
	for _, root := range roots {
		if _, err := os.Stat(root); err != nil {
			if !os.IsNotExist(err) {
				log.Warn().Err(err).Str("path", root).Msg("skipping root")
			} else {
				log.Debug().Str("path", root).Msg("root does not exist, skipping")
			}
			continue
		}
		submit(root)
	}

	// Close queue once all work is done so workers exit their range loops.
	go func() {
		wg.Wait()
		close(queue)
	}()

	// Worker pool: each worker processes one directory at a time.
	var workerWg sync.WaitGroup
	for range numWorkers {
		workerWg.Add(1)
		go func() {
			defer workerWg.Done()
			for dir := range queue {
				fdr.processDir(dir, glob, submit, results)
				wg.Done()
			}
		}()
	}

	// Close results after all workers finish so the collector goroutine exits.
	go func() {
		workerWg.Wait()
		close(results)
	}()

	var matches []string
	for r := range results {
		matches = append(matches, r)
	}

	return matches, nil
}

// processDir reads dir, sends any entries matching glob to results, and
// (if no match was found) submits subdirectories to submit for further
// traversal. It calls submit once per child, which increments wg; the
// caller is responsible for calling wg.Done() for dir itself.
func (fdr *Finder) processDir(dir, glob string, submit func(string), results chan<- string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Debug().Str("path", dir).Msg("directory disappeared during scan")
		} else {
			log.Warn().Err(err).Str("path", dir).Msg("failed to read directory")
		}
		return
	}

	// Single pass: collect directory matches and candidate subdirs simultaneously.
	// subdirs is only submitted if no match was found in this directory.
	var subdirs []string
	var foundMatch bool
	for _, entry := range entries {
		name := entry.Name()
		if fdr.isIgnored(name) {
			continue
		}
		matched, err := filepath.Match(glob, name)
		if err != nil {
			log.Warn().Err(err).Str("pattern", glob).Str("name", name).Msg("invalid match pattern")
			continue
		}
		if matched {
			results <- filepath.Join(dir, name)
			foundMatch = true
		} else if entry.IsDir() {
			subdirs = append(subdirs, name)
		}
	}

	if foundMatch {
		return
	}

	for _, name := range subdirs {
		submit(filepath.Join(dir, name))
	}
}
