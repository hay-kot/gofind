package gofind

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ParsePath(p string) string {
	// Check if path starts with ~ if it does, replace it with the user's home directory
	if strings.HasPrefix(p, "~") {
		homedir := Must(os.UserHomeDir())
		p = filepath.Join(homedir, strings.TrimPrefix(p, "~"))
	}

	return p
}

type App struct {
	verbose bool
}

func (a *App) LogCritical(args ...any) {
	fmt.Println(args...)
}

func (a *App) LogVerbose(args ...any) {
	if a.verbose {
		fmt.Println(args...)
	}
}

func (a *App) CacheAll() error {
	config := ReadDefaultConfig()
	config.CacheAll()

	return nil
}

func (a *App) Run(entry string) string {
	a.verbose = false

	if a.verbose {
		var startTime = time.Now()
		defer func() {
			fmt.Println("Execution time:", time.Since(startTime))
		}()
	}

	config := ReadDefaultConfig()
	// Parse Args

	useDefault := false

	if entry == "" {
		a.LogVerbose("No arguments provided using default argument")
		useDefault = true
	}

	cmd := config.Default

	if !useDefault {
		cmd = entry
	}

	search := config.Commands[cmd]

	// Check if cache is expired
	if config.Cache[cmd].IsExpired() {
		a.LogVerbose("Cache is expired, searching for results and rebuilding cache")
		config.Cache[cmd] = CacheEntry{
			Matches: search.Results(),
			Expires: time.Now().Add(time.Hour * 12),
		}
		config.Save()
	}

	matches := config.Cache[cmd].Matches

	a.LogVerbose(fmt.Sprintf("Found %d matches", len(matches)))

	filter := FzfFilter{}
	result := filter.Find(matches)

	return result.Path
}

type Match struct {
	Name string
	Path string
}

type SearchEntry struct {
	Root     string `json:"root"`
	MatchStr string `json:"match"`
}

func (se SearchEntry) Results() []Match {
	return SearchFor(se)
}

func SearchFor(search SearchEntry) []Match {
	var matches []Match
	var p = ParsePath(search.Root)

	var results = Must(Finder(p, search.MatchStr))

	if len(results) == 0 {
		panic("No results found")
	}

	for _, result := range results {
		match := filepath.Dir(result)
		name := filepath.Base(match)

		if name == "" {
			continue
		}

		matches = append(matches, Match{
			Name: name,
			Path: match,
		})
	}

	return matches
}
