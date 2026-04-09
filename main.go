package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"runtime/pprof"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"

	"github.com/hay-kot/gofind/internal/commands"
	"github.com/hay-kot/gofind/internal/paths"
)

var (
	// Build information. Populated at build-time via -ldflags flag.
	version = "dev"
	commit  = "HEAD"
	date    = "now"
)

func build() string {
	if version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok {
			version = info.Main.Version
			for _, s := range info.Settings {
				switch s.Key {
				case "vcs.revision":
					commit = s.Value
				case "vcs.time":
					date = s.Value
				}
			}
		}
	}

	short := commit
	if len(commit) > 7 {
		short = commit[:7]
	}

	return fmt.Sprintf("%s (%s) %s", version, short, date)
}

// setupLogger configures the global zerolog logger. Returns a closer that
// must be called when the process exits to flush and close any log file.
func setupLogger(level string, logFile string, noColor bool) (func(), error) {
	parsedLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		return func() {}, fmt.Errorf("failed to parse log level: %w", err)
	}

	var output io.Writer = zerolog.ConsoleWriter{Out: os.Stderr, NoColor: noColor}

	if logFile != "" {
		logDir := filepath.Dir(logFile)
		if err := os.MkdirAll(logDir, 0o755); err != nil {
			return func() {}, fmt.Errorf("failed to create log directory: %w", err)
		}

		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			return func() {}, fmt.Errorf("failed to open log file: %w", err)
		}

		output = io.MultiWriter(
			zerolog.ConsoleWriter{Out: os.Stderr, NoColor: noColor},
			file,
		)

		log.Logger = log.Output(output).Level(parsedLevel)
		return func() { _ = file.Close() }, nil
	}

	log.Logger = log.Output(output).Level(parsedLevel)
	return func() {}, nil
}

func main() {
	os.Exit(run())
}

func run() int {
	// Default logger for any output before setupLogger is called in Before.
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	var (
		logLevel   string
		noColor    bool
		logFile    string
		profileOut string
	)
	flags := &commands.Flags{}

	// logCloser is assigned by the Before hook once setupLogger succeeds.
	logCloser := func() {}
	defer func() { logCloser() }()

	// profileCloser stops CPU profiling if --profile was requested.
	profileCloser := func() {}
	defer func() { profileCloser() }()

	app := &cli.Command{
		Name:                  "gofind",
		Usage:                 "an interactive search for directories using the filepath.Match function",
		Version:               build(),
		EnableShellCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-level",
				Usage:       "log level (debug, info, warn, error, fatal, panic)",
				Sources:     cli.EnvVars("LOG_LEVEL"),
				Value:       "warn",
				Destination: &logLevel,
			},
			&cli.BoolFlag{
				Name:        "no-color",
				Usage:       "disable colored output",
				Sources:     cli.EnvVars("NO_COLOR"),
				Destination: &noColor,
			},
			&cli.StringFlag{
				Name:        "log-file",
				Usage:       "path to log file (optional)",
				Sources:     cli.EnvVars("LOG_FILE"),
				Destination: &logFile,
			},
			&cli.StringFlag{
				Name:        "config",
				Usage:       "specify configuration path",
				Sources:     cli.EnvVars("GOFIND_CONFIG"),
				DefaultText: paths.ConfigPath(""),
				Aliases:     []string{"c"},
				Destination: &flags.ConfigFile,
			},
			&cli.StringFlag{
				Name:        "profile",
				Usage:       "write CPU profile to file",
				Sources:     cli.EnvVars("GOFIND_PROFILE"),
				Destination: &profileOut,
			},
		},
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			closer, err := setupLogger(logLevel, logFile, noColor)
			if err != nil {
				return ctx, err
			}
			logCloser = closer

			if profileOut != "" {
				f, err := os.Create(profileOut)
				if err != nil {
					return ctx, fmt.Errorf("could not create profile file: %w", err)
				}
				if err := pprof.StartCPUProfile(f); err != nil {
					_ = f.Close()
					return ctx, fmt.Errorf("could not start CPU profile: %w", err)
				}
				profileCloser = func() {
					pprof.StopCPUProfile()
					_ = f.Close()
				}
			}

			return ctx, nil
		},
	}

	commands.NewCacheCmd(flags).Register(app)
	commands.NewFindCmd(flags).Register(app)
	commands.NewConfigCmd(flags).Register(app)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, os.Args); err != nil {
		colorRed := "\033[38;2;215;95;107m"
		colorGray := "\033[38;2;163;163;163m"
		colorReset := "\033[0m"
		if noColor {
			colorRed = ""
			colorGray = ""
			colorReset = ""
		}
		fmt.Fprintf(os.Stderr, "\n%s╭ Error%s\n%s│%s %s%s%s\n%s╵%s\n",
			colorRed, colorReset,
			colorRed, colorReset, colorGray, err.Error(), colorReset,
			colorRed, colorReset,
		)
		return 1
	}

	return 0
}
