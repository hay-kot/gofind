package gofind

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/hay-kot/yal"
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

func (gf *GoFind) CacheAll() error {
	cache := NewCache(gf.Conf.CacheDir)

	for key, entry := range gf.Conf.Commands {
		matches := gf.SearchFor(entry)
		cache.Set(key, matches)
		yal.Infof("Cached %v results for %s", len(matches), key)
	}
	return nil
}

func (gf *GoFind) Run(entry string) ([]Match, error) {
	useDefault := false

	if entry == "" {
		yal.Debug("no arguments provided using default argument")
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
			cached, _ = cache.Set(cmd, matches)
		} else {
			return []Match{}, err
		}
	}

	return cached.Matches, nil
}

func (gf *GoFind) SearchFor(search SearchEntry) []Match {
	var matches []Match

	paths := make([]string, len(search.Roots))
	for i, root := range search.Roots {
		paths[i] = ParsePath(root)
	}

	finder := Finder{
		MaxRecursion: 5,
		Ignore:       gf.Conf.Ignore,
	}

	var results = Must(finder.Find(paths, search.MatchStr))

	if len(results) == 0 {
		yal.Warnf("no results found for path %s", search.Roots)
		yal.Debugf("gf.SearchFor(Root=%s, MatchStr=%s) returned no results", search.Roots, search.MatchStr)
		return matches
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
