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

	// FlagErrorHandling ...
	FlagErrorHandling = flag.ContinueOnError

	errWarnNoArgs = errors.New("no more args")
	errWarnNoCmds = errors.New("no more cmds")
)

func nilWarnOrError(err error) error {
	if err == errWarnNoArgs || err == errWarnNoCmds {
		return nil
	}

	return err
}
