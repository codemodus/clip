package clip

import (
	"flag"
	"io/ioutil"
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
	for k := range subcmds.m {
		subcmds.m[k].pg = program
	}

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
			usage(c.pg, c.fs, subcmdsInfo(c.cs, ", "), err)
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
		return nil, nil, errCtrlNoArgs
	}

	nextArgs := args

	if c.fs != nil {
		if err := siftedParse(c.fs, args[1:], c.cs); err != nil {
			if isFlagHelpError(err) {
				c.no = true
			}
			return nil, nil, &FlagParseError{scp, err}
		}

		nextArgs = c.fs.Args()
		if len(nextArgs) == 0 {
			return nil, nil, errCtrlNoArgs
		}

		if c.cs == nil {
			if len(nextArgs) == 1 {
				return nil, nil, errCtrlNoCmds
			}

			return nil, nil, &BadCommandError{scp, c.fs.Arg(0)}
		}

		c.cs.cur = c.fs.Arg(0)
		if c.cs.cur == "" {
			return nil, nil, &EmptyCommandError{scp}
		}
	}

	nextCmd, ok := c.cs.m[c.cs.cur]
	if !ok {
		return nil, nil, &BadCommandError{scp, c.cs.cur}
	}

	return nextCmd, nextArgs, nil
}

func run(c *Command) (*Command, error) {
	scp := "run command"

	if c.fn != nil {
		if err := c.fn(); err != nil {
			return nil, err
		}
	}

	if c.no || c.cs == nil || len(c.cs.m) == 0 {
		return nil, errCtrlNoCmds
	}

	if c.cs.cur == "" {
		return nil, &EmptyCommandError{scp}
	}

	next, ok := c.cs.m[c.cs.cur]
	if !ok {
		return nil, &BadCommandError{scp, c.cs.cur}
	}

	return next, nil
}

func siftedParse(fs *flag.FlagSet, args []string, cs *CommandSet) error {
	out := fs.Output()
	fs.SetOutput(ioutil.Discard)
	defer fs.SetOutput(out)

	return fs.Parse(args)
}
