package clip

import "github.com/codemodus/clip/internal/clifsx"

// FlagSet represents a set of defined flags.
type FlagSet = clifsx.FlagSet

// NewFlagSet constructs a pointer to an instance of FlagSet.
func NewFlagSet(name string) *FlagSet {
	return clifsx.NewFlagSet(name)
}
