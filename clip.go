// Package clip provides simple command/subcommand structuring with a minimum
// of dependencies. The standard library flag.FlagSet type is leveraged often
// in this and related packages. If the FlagSet output is default, setting a
// help flag will result in usage output going to os.Stdout. Other parse errors
// will result in usage output going to os.Stderr as normal.
package clip

import (
	"fmt"

	"github.com/codemodus/clip/clipr"
	"github.com/codemodus/clip/internal/clifsx"
)

// Clip manages a CLI program which may or may not have a related function
// and/or subcommands. If a command does not have a related function and does
// have subcommands, it can be understood to be a command namespace.
type Clip struct {
	pg string
	fe bool
	tp bool
	fs *FlagSet
	fn func() error
	cs *CommandSet
}

// New constructs a pointer to an instance of Clip.
func New(program string, flags *FlagSet, subcmds *CommandSet) *Clip {
	setPrograms(program, subcmds)

	return &Clip{
		pg: program,
		tp: true,
		fs: flags,
		cs: subcmds,
	}
}

// Parse parses all provided args.
func (c *Clip) Parse(args []string) error {
	next, nextArgs, err := parse(c, args)
	if err != nil {
		if err = clipr.FilterControl(err); err != nil {
			err = clipr.NewUsageError(err, c.Usage, c.tp)
		}

		return err
	}

	return next.Parse(nextArgs)
}

// Usage recursively calls a command's FlagSet usages limited by depth where 0
// or less has the semantic meaning of full depth. If the causing error is of
// type flag.ErrHelp, the error is filtered and a nil error is returned. This
// enables convenient error returns for setting an application's exit code.
func (c *Clip) Usage(depth int, err error) error {
	return usage(c, 0, depth, err)
}

// UsageLongHelp calls a command's FlagSet usage. If usage is due to a help or
// h flag being set (as determined by the provided error), the command's
// FlagSet instances usages are called recursively.
func (c *Clip) UsageLongHelp(err error) error {
	if uerr, ok := AsUsageError(err); ok {
		depth := 1
		if uerr.IsHelp() {
			depth = 0
		}

		err = uerr.Usage(depth)
	}

	return err
}

// Run recursively runs all available *Command functions.
func (c *Clip) Run() error {
	next, err := run(c)
	if err != nil {
		return clipr.FilterControl(err)
	}

	return next.Run()
}

// HandlerFunc describes a function intended to be run as an endpoint of a
// called command.
type HandlerFunc func() error

// Command is an alias of Clip
type Command = Clip

// NewCommand constructs a pointer to an instance of Command.
func NewCommand(flags *FlagSet, fn HandlerFunc, subcmds *CommandSet) *Command {
	return &Command{
		pg: "unknown program",
		fs: flags,
		fn: fn,
		cs: subcmds,
	}
}

// NewCommandNamespace constructs a pointer to an instance of Command that is
// not intended to have a function associated with it.
func NewCommandNamespace(name string, subcmds *CommandSet) *Command {
	return &Command{
		pg: "unknown program namespace",
		fs: NewFlagSet(name),
		fn: nil,
		cs: subcmds,
	}
}

// CommandSet manages a currently scheduled command and available commands.
type CommandSet struct {
	cur string
	m   map[string]*Command
}

// NewCommandSet constructs a pointer to an instance of CommandSet.
func NewCommandSet(cmds ...*Command) *CommandSet {
	return &CommandSet{
		m: cmdsTable(cmds),
	}
}

// parse parses the provided *Command and returns the next requested *Command.
func parse(c *Command, args []string) (*Command, []string, error) {
	scp := "parse args"

	if args == nil || len(args) <= 1 {
		if c.fn == nil {
			return nil, nil, clipr.NewEmptyCommandError(scp)
		}

		return nil, nil, clipr.ErrCtrlNoArgs
	}

	nextArgs := args

	if c.fs != nil {
		if err := clifsx.Parse(c.fs, args[1:]); err != nil {
			if clipr.IsFlagHelp(err) {
				c.fe = true
			}
			return nil, nil, clipr.NewFlagParseError(scp, err)
		}

		nextArgs = c.fs.Args()
		if len(nextArgs) == 0 {
			return nil, nil, clipr.ErrCtrlNoArgs
		}

		if c.cs == nil {
			if len(nextArgs) == 1 {
				return nil, nil, clipr.ErrCtrlNoCmds
			}

			return nil, nil, clipr.NewBadCommandError(scp, c.fs.Arg(0))
		}

		c.cs.cur = c.fs.Arg(0)
		if c.cs.cur == "" {
			return nil, nil, clipr.NewEmptyCommandError(scp)
		}
	}

	nextCmd, ok := c.cs.m[c.cs.cur]
	if !ok {
		return nil, nil, clipr.NewBadCommandError(scp, c.cs.cur)
	}

	return nextCmd, nextArgs, nil
}

// usage recursively calls a command's FlagSet usages limited by depthToGo
// where 0 or less has the semantic meaning of full depth.
func usage(c *Command, depthGone, depthToGo int, err error) error {
	pg := c.pg
	if depthGone > 0 {
		pg = ""
	}

	rerr := clifsx.Usage(pg, depthGone, c.fs, subcmdsInfo(c.cs, ", "), err)

	if depthToGo == 1 || c.cs == nil {
		return rerr
	}

	for _, cmd := range c.cs.m {
		_ = usage(cmd, depthGone+1, depthToGo-1, err) //nolint
	}

	return rerr
}

// run runs the available *Command function and returns the next scheduled
// *Command.
func run(c *Command) (*Command, error) {
	scp := "run command"

	if c.fe {
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
		return nil, clipr.NewEmptyCommandError(scp)
	}

	next, ok := c.cs.m[c.cs.cur]
	if !ok {
		return nil, clipr.NewBadCommandError(scp, c.cs.cur)
	}

	return next, nil
}

// setPrograms recursively ensures all commands in a *CommandSet have the same
// program set.
func setPrograms(program string, subcmds *CommandSet) {
	if subcmds == nil || subcmds.m == nil {
		return
	}

	for k := range subcmds.m {
		subcmds.m[k].pg = program
		setPrograms(program, subcmds.m[k].cs)
	}
}

// cmdsTable converts a slice *Command to map for easy lookup.
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

// subcmdsInfo formats one level of commands within a *CommandSet.
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
