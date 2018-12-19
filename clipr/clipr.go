package clipr

import (
	"errors"
	"flag"
	"fmt"
)

var (
	// ErrCtrlNoArgs ...
	ErrCtrlNoArgs = errors.New("no more args")
	// ErrCtrlNoCmds ...
	ErrCtrlNoCmds = errors.New("no more cmds")
)

// EmptyCommandError ...
type EmptyCommandError struct {
	Scp string
}

// NewEmptyCommandError ...
func NewEmptyCommandError(scope string) *EmptyCommandError {
	return &EmptyCommandError{scope}
}

func (e *EmptyCommandError) Error() string {
	return fmt.Sprintf("%s: cannot use empty command", e.Scp)
}

// FlagParseError ...
type FlagParseError struct {
	Scp string
	Err error
}

// NewFlagParseError ...
func NewFlagParseError(scope string, err error) *FlagParseError {
	return &FlagParseError{scope, err}
}

func (e *FlagParseError) Error() string {
	return fmt.Sprintf("%s: %s", e.Scp, e.Err.Error())
}

// BadCommandError ...
type BadCommandError struct {
	Scp string
	Cmd string
}

// NewBadCommandError ...
func NewBadCommandError(scope, cmd string) *BadCommandError {
	return &BadCommandError{scope, cmd}
}

func (e *BadCommandError) Error() string {
	return fmt.Sprintf("%s: cannot find command: %s", e.Scp, e.Cmd)
}

// IsFlagHelpError ...
func IsFlagHelpError(err error) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(*FlagParseError); ok && e.Err == flag.ErrHelp {
		return true
	}

	return err == flag.ErrHelp
}

func isControlError(err error) bool {
	return err == ErrCtrlNoArgs || err == ErrCtrlNoCmds
}

// FilterControlError ...
func FilterControlError(err error) error {
	if isControlError(err) {
		return nil
	}

	return err
}
