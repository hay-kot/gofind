package tui

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hay-kot/gofind/gofind"
	"github.com/hay-kot/gofind/ui"
	"github.com/sahilm/fuzzy"
)

var (
	NoResults = errors.New("no results")
)

func FuzzyFinder(matches []gofind.Match) (string, error) {
	setup()

	ctrl := fuzzyFinderController{
		matches:  matches,
		selected: 0,
	}

	p := tea.NewProgram(
		newFuzzyFinderView(&ctrl),
		tea.WithMouseCellMotion(),
		tea.WithAltScreen(),
		tea.WithOutput(os.Stderr),
	)

	err := p.Start()

	selected := ctrl.Selected()

	if selected.Name == "" {
		return "", NoResults
	}

	return selected.Path, err
}

type fuzzyFinderController struct {
	index    []string
	matches  []gofind.Match
	selected int
	// key = filtered index, value = original index
	indexmap map[int]int
}

// Selected returns the active selection by the user, or any empty object
// if no selection has been made OR the active index is out of range.
func (c *fuzzyFinderController) Selected() gofind.Match {
	if c.indexmap == nil {
		if c.selected < 0 || c.selected >= len(c.matches) {
			return gofind.Match{}
		}
		return c.matches[c.selected]
	}

	idx, ok := c.indexmap[c.selected]
	if !ok {
		return gofind.Match{}
	}

	if idx < 0 || idx >= len(c.matches) {
		return gofind.Match{}
	}

	return c.matches[idx]
}

// Search returns a sorted list of matches uses a fuzzy search algorithm
func (c *fuzzyFinderController) Search(str string) []gofind.Match {
	if str == "" {
		return c.matches
	}

	c.indexmap = make(map[int]int)

	if c.index == nil {
		c.index = make([]string, len(c.matches))
		for i, repo := range c.matches {
			c.index[i] = repo.Name
		}
	}

	matches := fuzzy.Find(str, c.index)
	results := make([]gofind.Match, len(matches))
	for i, match := range matches {
		results[i] = c.matches[match.Index]
		c.indexmap[i] = match.Index
	}

	return results
}

type fuzzyFinderView struct {
	ctrl   *fuzzyFinderController
	search textinput.Model
	height int
	shift  int
	limit  int
}

func newFuzzyFinderView(ctrl *fuzzyFinderController) fuzzyFinderView {
	ti := textinput.New()
	ti.Focus()
	ti.Prompt = ui.AccentBlue("> ")
	ti.CharLimit = 256
	ti.Width = 80

	return fuzzyFinderView{
		search: ti,
		ctrl:   ctrl,
	}
}

func (m fuzzyFinderView) Init() tea.Cmd {
	return textinput.Blink
}

func (m fuzzyFinderView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			return m, tea.Quit
		}

		switch msg.String() {
		case "up":
			if m.ctrl.selected > 0 {
				m.ctrl.selected--

				if m.shift > 0 {
					m.shift--
				}
			}
		case "down":
			if m.ctrl.selected < len(m.ctrl.matches)-1 {
				m.ctrl.selected++

				if m.ctrl.selected >= len(m.ctrl.matches) {
					m.shift++
				}
			}
		}
	}

	m.search, cmd = m.search.Update(msg)
	return m, cmd
}

func (m fuzzyFinderView) View() string {
	results := m.ctrl.Search(m.search.Value())
	str := strings.Builder{}

	// Calculate the number of allowed_rows we can display
	m.limit = m.height - 12

	str.WriteString(
		m.search.View(),
	)

	str.WriteString(ui.Subtle(fmt.Sprintf("\n  %d/%d", len(results), len(m.ctrl.matches))) + "\n")
	str.WriteString(m.fmtMatches(results))

	return str.String()
}

func (m fuzzyFinderView) fmtMatches(repos []gofind.Match) string {
	longest := 0

	for _, repo := range repos {
		if len(repo.Name) > longest {
			longest = len(repo.Name)
		}
	}

	str := strings.Builder{}
	for i, repo := range repos {
		spaces := (longest + 5) - len(repo.Name)

		prefix := " "
		text := prefix + repo.Name + strings.Repeat(" ", spaces) + repo.Path
		if m.ctrl.selected == i {
			prefix = ui.HighlightRow(ui.AccentRed(">"))
			text = ui.HighlightRow(ui.Bold(text))
		}

		str.WriteString(prefix + text + "\n")
	}

	return str.String()
}
