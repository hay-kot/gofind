package commands

// Flags holds global flags shared across all commands.
type Flags struct {
	LogLevel   string
	NoColor    bool
	LogFile    string
	ConfigFile string
}
