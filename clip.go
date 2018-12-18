package clip

import (
	"flag"
)

// Clip ...
type Clip struct {
	fs *flag.FlagSet
	fn func() error
	cs *CommandSet
}

// New ...
func New(flags *flag.FlagSet, subCmds *CommandSet) *Clip {
	return &Clip{
		fs: flags,
		cs: subCmds,
	}
}

// Parse ...
func (c *Clip) Parse(args []string) error {
	next, args, err := parse(c, args)
	if err != nil {
		return nilWarnOrError(err)
	}

	return next.Parse(args)
}

// Run ...
func (c *Clip) Run() error {
	next, err := run(c)
	if err != nil {
		return nilWarnOrError(err)
	}

	return next.Run()
}

// CommandFunc ...
type CommandFunc func() error

// Command ...
type Command = Clip

// NewCommand ...
func NewCommand(flags *flag.FlagSet, fn CommandFunc, subCmds *CommandSet) *Command {
	return &Command{
		fs: flags,
		fn: fn,
		cs: subCmds,
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
	if args == nil || len(args) <= 1 {
		return nil, nil, errWarnNoArgs
	}

	nextArgs := args

	if c.fs != nil {
		if err := c.fs.Parse(args[1:]); err != nil {
			return nil, nil, ErrFlagParse
		}

		nextArgs = c.fs.Args()
		if len(nextArgs) == 0 {
			return nil, nil, errWarnNoArgs
		}

		if c.cs == nil {
			if len(nextArgs) == 1 {
				return nil, nil, errWarnNoCmds
			}

			return nil, nil, ErrBadCommand
		}

		c.cs.cur = c.fs.Arg(0)
		if c.cs.cur == "" {
			return nil, nil, ErrEmptyCommand
		}
	}

	nextCmd, ok := c.cs.m[c.cs.cur]
	if !ok {
		return nil, nil, ErrBadCommand
	}

	return nextCmd, nextArgs, nil
}

func run(c *Command) (*Command, error) {
	if c.fn != nil {
		if err := c.fn(); err != nil {
			return nil, err
		}
	}

	if c.cs == nil || len(c.cs.m) == 0 {
		return nil, errWarnNoCmds
	}

	if c.cs.cur == "" {
		return nil, ErrEmptyCommand
	}

	next, ok := c.cs.m[c.cs.cur]
	if !ok {
		return nil, ErrBadCommand
	}

	return next, nil
}
