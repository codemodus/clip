package clifs_test

import (
	"fmt"
	"os"
	"path"

	"github.com/codemodus/clip/clifs"
)

func Example() {
	// emulate cli command 'myapp -example=test -h'
	os.Args = []string{"myapp", "-example=test", "-h"}
	main()
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	var (
		program = path.Base(os.Args[0])
		verbose bool
		example string
	)

	fs := clifs.NewFlagSet("global")
	fs.BoolVar(&verbose, "v", verbose, "verbosity")
	fs.StringVar(&example, "example", example, "example")

	if err := clifs.Parse(fs, os.Args[1:]); err != nil {
		return clifs.Usage(program, fs, "", err)
	}

	fmt.Println(verbose, example)
	return nil
}
