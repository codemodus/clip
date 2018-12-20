package clifs

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/codemodus/clip/clipr"
)

var (
	// FlagErrorHandling ...
	FlagErrorHandling = flag.ContinueOnError
)

// Parse ...
func Parse(fs *flag.FlagSet, args []string) error {
	args = shiftCollision(args, fs.Name(), os.Args[0], path.Base(os.Args[0]))

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

	if clipr.IsFlagHelp(err) {
		err = nil

		if fs.Output() == os.Stderr {
			out := fs.Output()
			fs.SetOutput(os.Stdout)
			defer fs.SetOutput(out)
		}
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

func shiftCollision(args []string, ss ...string) []string {
	for _, s := range ss {
		if len(args) > 0 && args[0] == s {
			args = args[1:]
			break
		}
	}

	return args
}
