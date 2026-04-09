package gofind

import (
	"os"
	"path/filepath"
	"testing"
)

func homeDir(t testing.TB) string {
	t.Helper()
	h, err := os.UserHomeDir()
	if err != nil {
		t.Skipf("cannot determine home dir: %v", err)
	}
	return h
}

// BenchmarkFindGit benchmarks scanning the Code/repos and Code/pr directories
// for .git dirs, mirroring the default "repos" config entry.
func BenchmarkFindGit(b *testing.B) {
	home := homeDir(b)
	roots := []string{
		filepath.Join(home, "Code", "repos"),
		filepath.Join(home, "Code", "pr"),
	}

	var existing []string
	for _, r := range roots {
		if _, err := os.Stat(r); err == nil {
			existing = append(existing, r)
		}
	}
	if len(existing) == 0 {
		b.Skip("no Code/repos or Code/pr directories found")
	}

	fdr := &Finder{Ignore: []string{".git", "node_modules", "vendor"}}

	b.ResetTimer()
	for b.Loop() {
		_, err := fdr.Find(existing, ".git")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkFindGitSingle benchmarks scanning a single root directory for .git dirs.
func BenchmarkFindGitSingle(b *testing.B) {
	home := homeDir(b)
	root := filepath.Join(home, "Code", "repos")
	if _, err := os.Stat(root); err != nil {
		b.Skip("~/Code/repos not found")
	}

	fdr := &Finder{Ignore: []string{".git", "node_modules", "vendor"}}

	b.ResetTimer()
	for b.Loop() {
		_, err := fdr.Find([]string{root}, ".git")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkIgnorePatternMatch benchmarks the inner ignore-check loop.
func BenchmarkIgnorePatternMatch(b *testing.B) {
	ignore := []string{"node_modules", "vendor", ".git", "dist", "build", "target", "__pycache__"}
	name := "node_modules"

	b.ResetTimer()
	for b.Loop() {
		for _, pat := range ignore {
			matched, _ := filepath.Match(pat, name)
			if matched {
				break
			}
		}
	}
}
