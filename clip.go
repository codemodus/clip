package clip

import (
	"flag"
	"fmt"

	"github.com/codemodus/clip/clipr"
	"github.com/codemodus/clip/internal/clifsx"
)

var (
	// FlagErrorHandling ...
	FlagErrorHandling = flag.ContinueOnError
)

// Clip ...
type Clip struct {
	pg string
	fe bool
	tp bool
	fs *flag.FlagSet
	fn func() error
	cs *CommandSet
}

// New ...
func New(program string, flags *flag.FlagSet, subcmds *CommandSet) *Clip {
	setPrograms(program, subcmds)

	return &Clip{
		pg: program,
		tp: true,
		fs: flags,
		cs: subcmds,
	}
}

// Parse ...
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

// NoisyParse ...
func (c *Clip) NoisyParse(args []string) error {
	if err := c.Parse(args); err != nil {
		if uerr, ok := AsUsageError(err); ok {
			depth := 1
			if uerr.IsHelp() {
				depth = 0
			}

			err = uerr.Usage(depth)
		}

		return err
	}

	return nil
}

// Usage ...
func (c *Clip) Usage(depth int, err error) error {
	return usage(c, 0, depth, err)
}

// Run ...
func (c *Clip) Run() error {
	next, err := run(c)
	if err != nil {
		return clipr.FilterControl(err)
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

// UsageError ...
type UsageError interface {
	error
	Err() error
	Usage(depth int) error
	IsHelp() bool
	IsTop() bool
}

// AsUsageError ...
func AsUsageError(err error) (UsageError, bool) {
	uerr, ok := err.(UsageError)
	return uerr, ok
}

var (
	_ UsageError = (*clipr.UsageError)(nil)
)

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

func usage(c *Clip, depthGone, depthToGo int, err error) error {
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

func setPrograms(program string, subcmds *CommandSet) {
	if subcmds == nil || subcmds.m == nil {
		return
	}

	for k := range subcmds.m {
		subcmds.m[k].pg = program
		setPrograms(program, subcmds.m[k].cs)
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
