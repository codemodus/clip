package clip

import "github.com/codemodus/clip/clipr"

// UsageError ...
type UsageError interface {
	error
	Err() error
	Usage(depth int) error
	IsHelp() bool
	IsTop() bool
}

// AsUsageError ...
func AsUsageError(err error) (UsageError, bool) {
	uerr, ok := err.(UsageError)
	return uerr, ok
}

var (
	_ UsageError = (*clipr.UsageError)(nil)
)
