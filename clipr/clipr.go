// Package clipr provides error implementations and helpers related to the clip
// and clifs packages.
package clipr

import (
	"errors"
	"flag"
	"fmt"
)

type ctrlError error

var (
	// ErrCtrlNoArgs is a token signifying a control error due to lack of args.
	ErrCtrlNoArgs ctrlError = errors.New("no more args")
	// ErrCtrlNoCmds is a token signifying a control error due to lack of cmds.
	ErrCtrlNoCmds ctrlError = errors.New("no more cmds")
	// ErrHelp is a token signifying an error due to help request.
	ErrHelp = errors.New("help requested")
)

// NoBehaviorError manages no behavior error info.
type NoBehaviorError struct {
	Scp string
}

// NewNoBehaviorError constructs a pointer to an instance of NoBehaviorError.
func NewNoBehaviorError(scope string) *NoBehaviorError {
	return &NoBehaviorError{scope}
}

func (e *NoBehaviorError) Error() string {
	return fmt.Sprintf("%s: command has no defined behavior", e.Scp)
}

// EmptyCommandError manages empty command error info.
type EmptyCommandError struct {
	Scp string
}

// NewEmptyCommandError constructs a pointer to an instance of
// EmptyCommandError.
func NewEmptyCommandError(scope string) *EmptyCommandError {
	return &EmptyCommandError{scope}
}

func (e *EmptyCommandError) Error() string {
	return fmt.Sprintf("%s: cannot use empty command", e.Scp)
}

// FlagParseError manages flag parse error info.
type FlagParseError struct {
	Scp string
	Err error
}

// NewFlagParseError constructs a pointer to an instance of FlagParseError.
func NewFlagParseError(scope string, err error) *FlagParseError {
	return &FlagParseError{scope, err}
}

func (e *FlagParseError) Error() string {
	return fmt.Sprintf("%s: %s", e.Scp, e.Err.Error())
}

// BadCommandError manages bad command error info.
type BadCommandError struct {
	Scp string
	Cmd string
}

// NewBadCommandError constructs a pointer to an instance of BadCommandError.
func NewBadCommandError(scope, cmd string) *BadCommandError {
	return &BadCommandError{scope, cmd}
}

func (e *BadCommandError) Error() string {
	return fmt.Sprintf("%s: cannot find command: %s", e.Scp, e.Cmd)
}

// IsFlagHelp returns true if the error or underlying error is an instance of
// flag.ErrHelp.
func IsFlagHelp(err error) bool {
	if err == nil {
		return false
	}

	if ue, ok := err.(*UsageError); ok {
		err = ue.Err()
	}

	if fpe, ok := err.(*FlagParseError); ok && fpe.Err == flag.ErrHelp {
		return true
	}

	return err == flag.ErrHelp || err == ErrHelp
}

func isControl(err error) bool {
	return err == ErrCtrlNoArgs || err == ErrCtrlNoCmds
}

// FilterControl returns nil if the provided is a control error. Otherwise, the
// provided error is passed through.
func FilterControl(err error) error {
	if isControl(err) {
		return nil
	}

	return err
}

// UsageError manages usage error info.
type UsageError struct {
	err   error
	usage func(int, error) error
	top   bool
}

// NewUsageError constructs a pointer to an instance of UsageError. Top is to be
// understood as whether the provided error spawns from a type that is
// conceptually top-level.
func NewUsageError(err error, usageFunc func(int, error) error, top bool) *UsageError {
	return &UsageError{err, usageFunc, top}
}

func (e *UsageError) Error() string {
	return e.err.Error()
}

// Err returns the underlying error.
func (e *UsageError) Err() error {
	return e.err
}

// Usage calls the usage function of the type related to the underlying error.
func (e *UsageError) Usage(depth int) error {
	return e.usage(depth, e.err)
}

// IsHelp returns true if the underlying error or it's underlying error is an
// instance of flag.ErrHelp.
func (e *UsageError) IsHelp() bool {
	return IsFlagHelp(e.err)
}

// IsTop returns true if the type related to the underlying error is reported as
// being conceptually top-level.
func (e *UsageError) IsTop() bool {
	return e.top
}
