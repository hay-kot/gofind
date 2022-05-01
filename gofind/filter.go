package gofind

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type InteractiveFilter interface {
	Find(match []Match) Match
}

func (FzfFilter) formatSearch(repos []Match) string {
	longest := 0

	for _, repo := range repos {
		if len(repo.Name) > longest {
			longest = len(repo.Name)
		}
	}

	searchList := ""
	for _, repo := range repos {
		spaces := (longest + 5) - len(repo.Name)

		text := repo.Name + strings.Repeat(" ", spaces) + repo.Path

		searchList += text + "\n"
	}

	return searchList
}

type FzfFilter struct{}

func (f FzfFilter) Find(repos []Match) Match {
	var parseMatchStr = func(line string) (string, string) {
		strs := strings.Split(line, "    ")

		name := strings.TrimSpace(strs[0])

		path := strings.TrimSpace(strs[len(strs)-1])

		return name, path
	}

	searchList := f.formatSearch(repos)

	command := fmt.Sprintf("echo '%s' | fzf", searchList)

	// pipe list of repo names to fzf and get result
	cmd := exec.Command("bash", "-c", command)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	bts, err := cmd.Output()

	if err != nil {
		log.Fatal(err)
	}

	name := strings.TrimSpace(string(bts))

	name, path := parseMatchStr(name)

	for _, repo := range repos {
		if repo.Path == path {
			return repo
		}
	}

	panic("Could not find repo")
}
