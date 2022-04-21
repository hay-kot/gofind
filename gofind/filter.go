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
	var parseName = func(line string) string {
		return strings.TrimSpace(strings.Split(line, "    ")[0])
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

	name = parseName(name)

	for _, repo := range repos {

		if repo.Name == name {
			return repo
		}
	}

	panic("Could not find repo")
}
