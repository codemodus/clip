package clip

import "github.com/codemodus/clip/internal/clifsx"

// FlagSet ...
type FlagSet = clifsx.FlagSet

// NewFlagSet ...
func NewFlagSet(name string) *FlagSet {
	return clifsx.NewFlagSet(name)
}
