package clifs

import (
	"github.com/codemodus/clip/internal/clifsx"
)

// FlagSet ...
type FlagSet = clifsx.FlagSet

// NewFlagSet ...
func NewFlagSet(name string) *FlagSet {
	return clifsx.NewFlagSet(name)
}

// Parse ...
func Parse(fs *FlagSet, args []string) error {
	return clifsx.Parse(fs, args)
}

// Usage ...
func Usage(program string, fs *FlagSet, extra string, err error) error {
	return clifsx.Usage(program, 0, fs, extra, err)
}
