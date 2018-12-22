// Package clip provides simple command/subcommand structuring with a minimum
// of dependencies. The standard library flag.FlagSet type is leveraged often
// in this and related packages. If the FlagSet output is default, setting a
// help flag will result in usage output going to os.Stdout. Other parse errors
// will result in usage output going to os.Stderr as normal.
package clip

import (
	"strings"

	"github.com/codemodus/clip/clipr"
	"github.com/codemodus/clip/internal/clifsx"
)

// Clip manages a CLI program which may or may not have a related function
// and/or subcommands. If a command does not have a related function and does
// have subcommands, it can be understood to be a command namespace.
type Clip struct {
	pg string // program
	hf bool   // flag help
	tl bool   // top level
	fs *FlagSet
	fn func() error // command func
	cs *CommandSet
}

// New constructs a pointer to an instance of Clip.
func New(program string, fs *FlagSet, subcmds *CommandSet) *Clip {
	if fs == nil {
		fs = NewFlagSet(program)
	}

	setPrograms(program, subcmds)

	return &Clip{
		pg: program,
		tl: true,
		fs: fs,
		cs: subcmds,
	}
}

// Parse parses all provided args.
func (c *Clip) Parse(args []string) error {
	next, nextArgs, err := parse(c, args)
	if err != nil {
		if err = clipr.FilterControl(err); err != nil {
			err = clipr.NewUsageError(err, c.Usage, c.tl)
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
func NewCommand(fs *FlagSet, fn HandlerFunc, subcmds *CommandSet) *Command {
	program := "unknown program"

	if fs == nil {
		fs = NewFlagSet(program)
	}

	return &Command{
		pg: program,
		fs: fs,
		fn: fn,
		cs: subcmds,
	}
}

// NewCommandNamespace constructs a pointer to an instance of Command that is
// not intended to have a function associated with it.
func NewCommandNamespace(name string, subcmds *CommandSet) *Command {
	return NewCommand(NewFlagSet(name), nil, subcmds)
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

	if isEmptyArgs(args) || len(args) == 1 {
		if c.fn == nil {
			return nil, nil, clipr.NewNoBehaviorError(scp)
		}

		return nil, nil, clipr.ErrCtrlNoArgs
	}

	if err := clifsx.Parse(c.fs, args[1:]); err != nil {
		c.hf = clipr.IsFlagHelp(err)
		return nil, nil, clipr.NewFlagParseError(scp, err)
	}

	return nextToParse(scp, c)
}

// nextToParse returns the next requested *Command.
func nextToParse(scp string, c *Command) (*Command, []string, error) {
	nextArgs := c.fs.Args()
	if isEmptyArgs(nextArgs) {
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

	cmdsInfo := commandSetInfo(c.cs, ", ", "Available commands - ")
	rerr := clifsx.Usage(pg, depthGone, c.fs, cmdsInfo, err)

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

	if c.hf {
		return nil, clipr.ErrHelp
	}

	if c.fn != nil {
		if err := c.fn(); err != nil {
			return nil, err
		}
	}

	if isEmptyCommandSet(c.cs) {
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
func setPrograms(program string, cs *CommandSet) {
	if isEmptyCommandSet(cs) {
		return
	}

	for k := range cs.m {
		cs.m[k].pg = program
		setPrograms(program, cs.m[k].cs)
	}
}

// cmdsTable converts a slice *Command to map for easy lookup.
func cmdsTable(cmds []*Command) map[string]*Command {
	m := make(map[string]*Command)

	if isEmptyCmds(cmds) {
		return m
	}

	for _, c := range cmds {
		if isNamedFlagSet(c.fs) {
			m[c.fs.Name()] = c
		}
	}

	return m
}

// commandSetInfo formats one level of commands within a *CommandSet.
func commandSetInfo(cs *CommandSet, sep string, prefix string) string {
	if isEmptyCommandSet(cs) {
		return ""
	}

	var b strings.Builder
	_, _ = b.WriteString(prefix) //nolint

	var sp string
	for k := range cs.m {
		_, _ = b.WriteString(sp + k) //nolint
		sp = sep
	}

	return b.String()
}

func isEmptyArgs(args []string) bool {
	return args == nil || len(args) == 0
}

func isEmptyCommandSet(cs *CommandSet) bool {
	return cs == nil || isEmptyCmdsTable(cs.m)
}

func isEmptyCmds(cmds []*Command) bool {
	return cmds == nil || len(cmds) == 0
}

func isNamedFlagSet(fs *FlagSet) bool {
	return fs != nil && fs.Name() != ""
}

func isEmptyCmdsTable(m map[string]*Command) bool {
	return m == nil || len(m) == 0
}
