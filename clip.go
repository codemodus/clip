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
	if args == nil || len(args) <= 1 {
		return nil
	}

	nextArgs := args

	if c.fs != nil {
		if err := c.fs.Parse(args[1:]); err != nil {
			return ErrFlagParse
		}

		nextArgs = c.fs.Args()
		if len(nextArgs) == 0 {
			return nil
		}

		if c.cs == nil {
			return ErrBadCommand
		}

		c.cs.cur = c.fs.Arg(0)

		if c.cs.cur == "" {
			return ErrEmptyCommand
		}
	}

	cc, ok := c.cs.m[c.cs.cur]
	if !ok {
		return ErrBadCommand
	}

	return cc.Parse(nextArgs)
}

// Run ...
func (c *Clip) Run() error {
	if c.fn != nil {
		if err := c.fn(); err != nil {
			return err
		}
	}

	if c.cs == nil || len(c.cs.m) == 0 {
		// TODO: check if required, parse user input if needed, then run
		return nil
	}

	if c.cs.cur == "" {
		return ErrEmptyCommand
	}

	next, ok := c.cs.m[c.cs.cur]
	if !ok {
		return ErrBadCommand
	}

	return next.Run()
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
