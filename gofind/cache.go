package gofind

import "time"

type CacheEntry struct {
	Matches []Match   `json:"matches"`
	Expires time.Time `json:"expires"`
}

func (ce CacheEntry) IsExpired() bool {
	return time.Now().After(ce.Expires)
}
