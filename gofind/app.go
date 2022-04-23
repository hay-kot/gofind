package gofind

import (
	"os"
	"path/filepath"
	"strings"
)

func ParsePath(p string) string {
	// Check if path starts with ~ if it does, replace it with the user's home directory
	if strings.HasPrefix(p, "~") {
		homedir := Must(os.UserHomeDir())
		p = filepath.Join(homedir, strings.TrimPrefix(p, "~"))
	}

	return p
}

type GoFind struct {
	Verbose bool
	Conf    Config
}

func (a *GoFind) CacheAll() error {
	return nil
}

func (gf *GoFind) Run(entry string) (string, error) {
	useDefault := false

	if entry == "" {
		gf.LogVerbose("No arguments provided using default argument")
		useDefault = true
	}

	cmd := gf.Conf.Default

	if !useDefault {
		cmd = entry
	}

	// Preload cache if exists
	cache := NewCache(gf.Conf.CacheDir)
	cached, err := cache.Find(cmd)

	if err != nil {
		if err == ErrCacheNotFound {
			matches := gf.SearchFor(gf.Conf.Commands[cmd])
			cached, err = cache.Set(cmd, matches)
		} else {
			return "", err
		}
	}

	filter := FzfFilter{}
	result := filter.Find(cached.Matches)

	return result.Path, nil
}

func (gf *GoFind) SearchFor(search SearchEntry) []Match {
	var matches []Match
	var p = ParsePath(search.Root)

	finder := Finder{
		MaxRecursion: 5,
		Ignore:       gf.Conf.Ignore,
	}

	var results = Must(finder.Find(p, search.MatchStr))

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

type Match struct {
	Name string
	Path string
}

type SearchEntry struct {
	Root     string `json:"root"`
	MatchStr string `json:"match"`
}
