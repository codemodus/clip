package clip

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

// Clip ...
type Clip struct {
	pg string
	no bool
	fs *flag.FlagSet
	fn func() error
	cs *CommandSet
}

// New ...
func New(program string, flags *flag.FlagSet, subcmds *CommandSet) *Clip {
	setPrograms(program, subcmds)

	return &Clip{
		pg: program,
		fs: flags,
		cs: subcmds,
	}
}

// Parse ...
func (c *Clip) Parse(args []string) error {
	next, nextArgs, err := parse(c, args)
	if err != nil {
		usg := func() {
			Usage(c.pg, c.fs, subcmdsInfo(c.cs, ", "), err)
		}

		return filteredError(usg, err)
	}

	return next.Parse(nextArgs)
}

// Run ...
func (c *Clip) Run() error {
	next, err := run(c)
	if err != nil {
		return filteredError(nil, err)
	}

	return next.Run()
}

// HandlerFunc ...
type HandlerFunc func() error

// Command ...
type Command = Clip

// NewCommand ...
func NewCommand(flags *flag.FlagSet, fn HandlerFunc, subcmds *CommandSet) *Command {
	return &Command{
		pg: "unknown program",
		fs: flags,
		fn: fn,
		cs: subcmds,
	}
}

// NewCommandNamespace ...
func NewCommandNamespace(name string, subcmds *CommandSet) *Command {
	return &Command{
		pg: "unknown program namespace",
		fs: flag.NewFlagSet(name, FlagErrorHandling),
		fn: nil,
		cs: subcmds,
	}
}

// CommandSet ...
type CommandSet struct {
	cur string
	m   map[string]*Command
}

// NewCommandSet ...
func NewCommandSet(cmds ...*Command) *CommandSet {
	return &CommandSet{
		m: cmdsTable(cmds),
	}
}

func cmdsTable(cmds []*Command) map[string]*Command {
	m := make(map[string]*Command)

	if cmds == nil {
		return m
	}

	for _, c := range cmds {
		if c.fs != nil && c.fs.Name() != "" {
			m[c.fs.Name()] = c
		}
	}

	return m
}

func parse(c *Command, args []string) (*Command, []string, error) {
	scp := "parse args"

	if args == nil || len(args) <= 1 {
		if c.fn == nil {
			return nil, nil, &clipr.EmptyCommandError{scp}
		}

		return nil, nil, clipr.ErrCtrlNoArgs
	}

	nextArgs := args

	if c.fs != nil {
		if err := Parse(c.fs, args[1:]); err != nil {
			if clipr.IsFlagHelpError(err) {
				c.no = true
			}
			return nil, nil, &clipr.FlagParseError{scp, err}
		}

		nextArgs = c.fs.Args()
		if len(nextArgs) == 0 {
			return nil, nil, clipr.ErrCtrlNoArgs
		}

		if c.cs == nil {
			if len(nextArgs) == 1 {
				return nil, nil, clipr.ErrCtrlNoCmds
			}

			return nil, nil, &clipr.BadCommandError{scp, c.fs.Arg(0)}
		}

		c.cs.cur = c.fs.Arg(0)
		if c.cs.cur == "" {
			return nil, nil, &clipr.EmptyCommandError{scp}
		}
	}

	nextCmd, ok := c.cs.m[c.cs.cur]
	if !ok {
		return nil, nil, &clipr.BadCommandError{scp, c.cs.cur}
	}

	return nextCmd, nextArgs, nil
}

func run(c *Command) (*Command, error) {
	scp := "run command"

	if c.no {
		return nil, clipr.ErrCtrlNoCmds
	}

	if c.fn != nil {
		if err := c.fn(); err != nil {
			return nil, err
		}
	}

	if c.cs == nil || len(c.cs.m) == 0 {
		return nil, clipr.ErrCtrlNoCmds
	}

	if c.cs.cur == "" {
		return nil, &clipr.EmptyCommandError{scp}
	}

	next, ok := c.cs.m[c.cs.cur]
	if !ok {
		return nil, &clipr.BadCommandError{scp, c.cs.cur}
	}

	return next, nil
}

// Parse ...
func Parse(fs *flag.FlagSet, args []string) error {
	out := fs.Output()
	fs.SetOutput(ioutil.Discard)
	defer fs.SetOutput(out)

	return fs.Parse(args)
}

// Usage ...
func Usage(program string, fs *flag.FlagSet, extra string, err error) {
	if clipr.IsFlagHelpError(err) && fs.Output() == os.Stderr {
		out := fs.Output()
		fs.SetOutput(os.Stdout)
		defer fs.SetOutput(out)
	}

	if program != "" && program != fs.Name() {
		fmt.Fprintf(fs.Output(), "%s:\n", program)
	}

	fs.Usage()

	if extra != "" {
		fmt.Fprintln(fs.Output(), extra)
	}
}

func setPrograms(program string, subcmds *CommandSet) {
	if subcmds == nil || subcmds.m == nil {
		return
	}

	for k := range subcmds.m {
		subcmds.m[k].pg = program
		setPrograms(program, subcmds.m[k].cs)
	}
}

func subcmdsInfo(cs *CommandSet, sep string) string {
	var s string

	if cs == nil || len(cs.m) == 0 {
		return s
	}

	for k := range cs.m {
		s += k + sep
	}

	if len(s) > len(sep)-1 {
		s = s[:len(s)-len(sep)]
	}

	return fmt.Sprintf("Available commands - %s", s)
}

func filteredError(usage func(), err error) error {
	if clipr.IsControlError(err) {
		return nil
	}

	if usage != nil {
		usage()
	}

	if clipr.IsFlagHelpError(err) {
		return nil
	}

	return err
}
