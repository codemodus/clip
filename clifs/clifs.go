// Package clifs wraps flag.FlagSet to improve the flexibility of usage output.
// If the FlagSet output is default, setting a help flag will result in usage
// output going to os.Stdout. Other parse errors will result in usage output
// going to os.Stderr as normal.
package clifs

import (
	"github.com/codemodus/clip/internal/clifsx"
)

// FlagSet represents a set of defined flags.
type FlagSet = clifsx.FlagSet

// NewFlagSet constructs a pointer to an instance of FlagSet.
func NewFlagSet(name string) *FlagSet {
	return clifsx.NewFlagSet(name)
}

// Parse calls the FlagSet parse function after silencing usage output. An
// attempt is made at removing accidentally included command names.
func Parse(fs *FlagSet, args []string) error {
	return clifsx.Parse(fs, args)
}

// Usage prints the program name, FlagSet usage, and extra line. If the causing
// error is of type flag.ErrHelp, the error is filtered and a nil error is
// returned. This enables convenient error returns for setting an application's
// exit code.
func Usage(program string, fs *FlagSet, extra string, err error) error {
	return clifsx.Usage(program, 0, fs, extra, err)
}
