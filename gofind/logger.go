package gofind

import "fmt"

func (a *GoFind) LogCritical(args ...any) {
	fmt.Println(args...)
}

func (a *GoFind) LogVerbose(args ...any) {
	if a.Verbose {
		fmt.Println(args...)
	}
}
