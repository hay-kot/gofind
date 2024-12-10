package gofind

import "github.com/rs/zerolog/log"

func Must[T any](v T, err error) T {
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}
	return v
}

func NoErr(err error) {
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}
}
