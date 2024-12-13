package ui

import "strings"

// Table takes in a 2D array of strings and returns a string of a table.
// including the header. It returns an evenly spaced table for neatly printing
// tables with a consistent and simple look.
func Table(rows [][]string) string {
	table := strings.Builder{}
	cols := len(rows[0])

	// Find longest string in each column
	max := make([]int, cols)
	for _, row := range rows {
		for i, s := range row {
			if len(s) > max[i] {
				max[i] = len(s)
			}
		}
	}

	for i, row := range rows {
		for j, s := range row {
			// Pad the string with spaces to the max length of the column + 4
			spaces := strings.Repeat(" ", max[j]-len(s)+4)

			if i == 0 {
				table.WriteString(Bold(s + spaces))
			} else {
				table.WriteString(s + spaces)
			}
		}
		table.WriteString("\n")
	}

	return table.String()
}
