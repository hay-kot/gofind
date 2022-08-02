package gofind

import "github.com/hay-kot/yal"

func Must[T any](v T, err error) T {
	if err != nil {
		yal.Fatal(err.Error())
	}
	return v
}

func NoErr(err error) {
	if err != nil {
		yal.Fatal(err.Error())
	}
}
