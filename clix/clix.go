package clix

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/codemodus/clip/clipr"
)

var (
	// FlagErrorHandling ...
	FlagErrorHandling = flag.ContinueOnError
)

// Parse ...
func Parse(fs *flag.FlagSet, args []string) error {
	out := fs.Output()
	fs.SetOutput(ioutil.Discard)
	defer fs.SetOutput(out)

	return fs.Parse(args)
}

// Usage ...
func Usage(program string, fs *flag.FlagSet, extra string, err error) error {
	if err == nil {
		return nil
	}

	if IsFlagHelpError(err) && fs.Output() == os.Stderr {
		out := fs.Output()
		fs.SetOutput(os.Stdout)
		defer fs.SetOutput(out)

		err = nil
	}

	if program != "" && program != fs.Name() {
		fmt.Fprintf(fs.Output(), "%s:\n", program)
	}

	fs.Usage()

	if extra != "" {
		fmt.Fprintln(fs.Output(), extra)
	}

	return err
}

// IsFlagHelpError ...
func IsFlagHelpError(err error) bool {
	return clipr.IsFlagHelpError(err)
}
