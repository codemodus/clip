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
func New(flags *flag.FlagSet, cmds *CommandSet) *Clip {
	return &Clip{
		fs: flags,
		cs: cmds,
	}
}

// Parse ...
func (c *Clip) Parse(args []string) error {
	next, args, err := parse(c, args)
	if err != nil {
		if err == errWarnNoArgs {
			return nil
		}

		return err
	}

	return next.Parse(args)
}

func parse(c *Clip, args []string) (*Command, []string, error) {
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

// Run ...
func (c *Clip) Run() error {
	next, err := run(c)
	if err != nil {
		if err == errWarnNoCmds {
			return nil
		}

		return err
	}

	return next.Run()
}

func run(c *Clip) (*Clip, error) {
	if c.fn != nil {
		if err := c.fn(); err != nil {
			return nil, err
		}
	}

	if c.cs == nil || len(c.cs.m) == 0 {
		// TODO: check if required, parse user input if needed, then run
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

// Command ...
type Command = Clip

// NewCommand ...
func NewCommand(flags *flag.FlagSet, fn func() error, cmds *CommandSet) *Command {
	return &Command{
		fs: flags,
		fn: fn,
		cs: cmds,
	}
}

// CommandSet ...
type CommandSet struct {
	req bool
	cur string
	m   map[string]*Command
}

// NewCommandSet ...
func NewCommandSet(required bool, cmds ...*Command) *CommandSet {
	m := make(map[string]*Command)

	for _, c := range cmds {
		if c.fs != nil && c.fs.Name() != "" {
			m[c.fs.Name()] = c
		}
	}

	return &CommandSet{
		req: required,
		m:   m,
	}
}
