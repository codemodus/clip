package clip_test

import (
	"flag"
	"fmt"
	"os"

	"github.com/codemodus/clip"
)

func Example() {
	var (
		globalCnf = newGlobalConf()
		printCnf  = newPrintConf("print")
		otherCnf  = newOtherConf("other")
	)

	cs := clip.NewCommandSet(true,
		clip.NewCommand(printCnf.flagSet, runPrintFunc(printCnf, globalCnf), nil),
		clip.NewCommand(otherCnf.flagSet, runOtherFunc(otherCnf), nil),
	)
	app := clip.New(globalCnf.flagSet, cs)

	// emulate cli command 'myapp -v print -msg=hello, world'
	os.Args = []string{"myapp", "-v", "print", "-msg=hello, world"}

	if err := app.Parse(os.Args); err != nil {
		// ...
	}

	if err := app.Run(); err != nil {
		// ...
	}

	// Output:
	// hello, world (global verbosity is enabled)
}

func ExampleCommandFunc() {
	var runPrint clip.CommandFunc
	var runAdvPrint clip.CommandFunc

	runPrint = func() error {
		_, err := fmt.Println("hello, example")
		return err
	}

	runAdvPrintFunc := func(msg string, verbosity bool) func() error {
		return func() error {
			_, err := fmt.Printf("%s (verbosity = %t)\n", msg, verbosity)
			return err
		}
	}
	runAdvPrint = runAdvPrintFunc("hello, again", true)

	runPrint()
	runAdvPrint()

	// Output:
	// hello, example
	// hello, again (verbosity = true)
}

type globalConf struct {
	flagSet *flag.FlagSet
	verbose bool
}

func newGlobalConf() *globalConf {
	c := globalConf{
		flagSet: flag.NewFlagSet("global", clip.FlagErrorHandling),
	}

	c.flagSet.BoolVar(&c.verbose, "v", c.verbose, "enable verbosity")

	return &c
}

type printConf struct {
	flagSet *flag.FlagSet
	msg     string
}

func newPrintConf(name string) *printConf {
	c := printConf{
		flagSet: flag.NewFlagSet(name, clip.FlagErrorHandling),
		msg:     "default message",
	}

	c.flagSet.StringVar(&c.msg, "msg", c.msg, "message to print")

	return &c
}

type otherConf struct {
	flagSet *flag.FlagSet
	file    string
}

func newOtherConf(name string) *otherConf {
	c := otherConf{
		flagSet: flag.NewFlagSet(name, clip.FlagErrorHandling),
		file:    "test_data",
	}

	return &c
}

func runPrintFunc(cnf *printConf, gCnf *globalConf) func() error {
	return func() error {
		stts := "disabled"
		if gCnf.verbose {
			stts = "enabled"
		}

		_, err := fmt.Printf("%s (global verbosity is %s)\n", cnf.msg, stts)
		return err

	}
}

func runOtherFunc(cnf *otherConf) func() error {
	return func() error {
		return nil
	}
}
