package clip

import (
	"errors"
	"flag"
	"os"
)

var (
	// ErrFlagParse ...
	ErrFlagParse = errors.New("failed to parse flags")

	// ErrBadCommand ...
	ErrBadCommand = errors.New("cannot find command")

	// FlagErrorHandling ...
	FlagErrorHandling = flag.ContinueOnError
)

// Clip ...
type Clip struct {
	fs *flag.FlagSet
	fn func() error
	cs *CommandSet
}

// New ...
func New(flags *flag.FlagSet, cmds *CommandSet) *Clip {
	return &Clip{fs: flags, cs: cmds}
}

// Parse ...
func (c *Clip) Parse(flags []string) error {
	if len(flags) <= 1 {
		return nil
	}

	if err := c.fs.Parse(flags[1:]); err != nil {
		return ErrFlagParse
	}

	if len(c.fs.Args()) == 0 {
		return nil
	}

	c.cs.cc = c.fs.Arg(0)

	cc, ok := c.cs.cs[c.cs.cc]
	if !ok {
		return ErrBadCommand
	}

	return cc.Parse(nextArgs(os.Args, c.cs.cc))
}

// Run ...
func (c *Clip) Run() error {
	if c.fn != nil {
		if err := c.fn(); err != nil {
			return err
		}
	}

	if c.cs == nil || len(c.cs.cs) == 0 {
		// TODO: check if required, parse user input if needed, then run
		return nil
	}

	next, ok := c.cs.cs[c.cs.cc]
	if !ok {
		// TODO: this probably can't happen, verify and tighten accordingly
		return nil
	}
	return next.Run()
}

// Command ...
type Command = Clip

// NewCommand ...
func NewCommand(flags *flag.FlagSet, fn func() error, cmds *CommandSet) *Command {
	return &Command{fs: flags, fn: fn, cs: cmds}
}

// CommandSet ...
type CommandSet struct {
	req bool
	cc  string
	cs  map[string]*Command
}

// NewCommandSet ...
func NewCommandSet(required bool, cmds ...*Command) *CommandSet {
	cs := make(map[string]*Command)
	for _, c := range cmds {
		cs[c.fs.Name()] = c
	}

	return &CommandSet{req: required, cs: cs}
}

func nextArgs(vals []string, val string) []string {
	for k, v := range vals {
		if v == val {
			return vals[k:]
		}
	}

	return vals
}
