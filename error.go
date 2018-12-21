package clip

import "github.com/codemodus/clip/clipr"

// UsageError manages parse error related contextual information.
type UsageError interface {
	error
	Err() error
	Usage(depth int) error
	IsHelp() bool
	IsTop() bool
}

// AsUsageError is a convenience function for asserting that an error is a
// UsageError.
func AsUsageError(err error) (UsageError, bool) {
	uerr, ok := err.(UsageError)
	return uerr, ok
}

var (
	_ UsageError = (*clipr.UsageError)(nil)
)
