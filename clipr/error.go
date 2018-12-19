package clipr

import (
	"errors"
	"flag"
	"fmt"
)

var (
	ErrCtrlNoArgs = errors.New("no more args")
	ErrCtrlNoCmds = errors.New("no more cmds")
)

// EmptyCommandError ...
type EmptyCommandError struct {
	Scp string
}

func (e *EmptyCommandError) Error() string {
	return fmt.Sprintf("%s: cannot use empty command", e.Scp)
}

// FlagParseError ...
type FlagParseError struct {
	Scp string
	Err error
}

func (e *FlagParseError) Error() string {
	return fmt.Sprintf("%s: %s", e.Scp, e.Err.Error())
}

// BadCommandError ...
type BadCommandError struct {
	Scp string
	Cmd string
}

func (e *BadCommandError) Error() string {
	return fmt.Sprintf("%s: cannot find command: %s", e.Scp, e.Cmd)
}

func IsFlagHelpError(err error) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(*FlagParseError); ok && e.Err == flag.ErrHelp {
		return true
	}

	return err == flag.ErrHelp
}

func IsControlError(err error) bool {
	return err == ErrCtrlNoArgs || err == ErrCtrlNoCmds
}
