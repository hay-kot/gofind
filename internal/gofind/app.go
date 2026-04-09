package gofind

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/hay-kot/gofind/internal/config"
	"github.com/rs/zerolog/log"
)

func ParsePath(p string) (string, error) {
	// Check if path starts with ~ if it does, replace it with the user's home directory
	if strings.HasPrefix(p, "~") {
		homedir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		p = filepath.Join(homedir, strings.TrimPrefix(p, "~"))
	}

	return p, nil
}

type GoFind struct {
	Verbose bool
	NoCache bool
	Conf    *config.Config
}

func (gf *GoFind) CacheAll() error {
	cache, err := NewCache(gf.Conf.CacheDir)
	if err != nil {
		return err
	}

	for key, entry := range gf.Conf.Commands {
		matches := gf.SearchFor(entry)
		_, err := cache.Set(key, matches)
		if err != nil {
			return err
		}
		log.Info().Str("key", key).Int("matches", len(matches)).Msg("cached results")
	}
	return nil
}

func (gf *GoFind) Run(entry string) ([]Match, error) {
	cmd := entry
	if entry == "" {
		log.Debug().Msg("no arguments provided using default")
		cmd = gf.Conf.Default
	}

	if gf.NoCache {
		return gf.SearchFor(gf.Conf.Commands[cmd]), nil
	}

	// Preload cache if exists
	cache, err := NewCache(gf.Conf.CacheDir)
	if err != nil {
		return []Match{}, err
	}
	cached, err := cache.Find(cmd)
	if err != nil {
		if errors.Is(err, ErrCacheNotFound) {
			matches := gf.SearchFor(gf.Conf.Commands[cmd])
			cached, _ = cache.Set(cmd, matches)
		} else {
			return []Match{}, err
		}
	}

	return cached.Matches, nil
}

func (gf *GoFind) SearchFor(search config.SearchEntry) []Match {
	var matches []Match

	paths := make([]string, 0, len(search.Roots))
	for _, root := range search.Roots {
		parsedPath, err := ParsePath(root)
		if err != nil {
			log.Warn().Err(err).Str("root", root).Msg("failed to parse path")
			continue
		}
		paths = append(paths, parsedPath)
	}

	finder := Finder{
		Ignore: gf.Conf.Ignore,
	}

	results, err := finder.Find(paths, search.MatchStr)
	if err != nil {
		log.Warn().Err(err).Msg("finder.Find failed")
		return matches
	}

	if len(results) == 0 {
		log.Warn().Strs("path", search.Roots).Msg("no results found")
		log.Debug().Msgf("gf.SearchFor(Root=%s, MatchStr=%s) returned no results", search.Roots, search.MatchStr)
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

	slices.SortFunc(matches, func(a, b Match) int {
		return strings.Compare(a.Name, b.Name)
	})
	return matches
}
