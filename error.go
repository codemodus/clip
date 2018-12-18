package clip

import (
	"errors"
	"flag"
)

var (
	// ErrFlagParse ...
	ErrFlagParse = errors.New("cannot parse flags")
	// ErrBadCommand ...
	ErrBadCommand = errors.New("cannot find command")
	// ErrEmptyCommand ...
	ErrEmptyCommand = errors.New("cannot use empty command")
)

var (
	// FlagErrorHandling ...
	FlagErrorHandling = flag.ContinueOnError
)

var (
	errWarnNoArgs = errors.New("no more args")
	errWarnNoCmds = errors.New("no more cmds")
)
