package clifs

import (
	"flag"

	"github.com/codemodus/clip/internal/clifsx"
)

var (
	// FlagErrorHandling ...
	FlagErrorHandling = flag.ContinueOnError
)

// Parse ...
func Parse(fs *flag.FlagSet, args []string) error {
	return clifsx.Parse(fs, args)
}

// Usage ...
func Usage(program string, fs *flag.FlagSet, extra string, err error) error {
	return clifsx.Usage(program, 0, fs, extra, err)
}
