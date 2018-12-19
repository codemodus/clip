package clip

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var (
	// FlagErrorHandling ...
	FlagErrorHandling = flag.ContinueOnError

	errCtrlNoArgs = errors.New("no more args")
	errCtrlNoCmds = errors.New("no more cmds")
)

// EmptyCommandError ...
type EmptyCommandError struct {
	scp string
}

func (e *EmptyCommandError) Error() string {
	return fmt.Sprintf("%s: cannot use empty command", e.scp)
}

// FlagParseError ...
type FlagParseError struct {
	scp string
	err error
}

func (e *FlagParseError) Error() string {
	return fmt.Sprintf("%s: %s", e.scp, e.err.Error())
}

// BadCommandError ...
type BadCommandError struct {
	scp string
	cmd string
}

func (e *BadCommandError) Error() string {
	return fmt.Sprintf("%s: cannot find command: %s", e.scp, e.cmd)
}

func isFlagHelpError(err error) bool {
	if err == nil {
		return false
	}

	if e, ok := err.(*FlagParseError); ok && e.err == flag.ErrHelp {
		return true
	}

	return err == flag.ErrHelp
}

func isControlError(err error) bool {
	return err == errCtrlNoArgs || err == errCtrlNoCmds
}

func filteredError(usage func(), err error) error {
	if isControlError(err) {
		return nil
	}

	if usage != nil {
		usage()
	}

	if isFlagHelpError(err) {
		return nil
	}

	return err
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

func usage(program string, fs *flag.FlagSet, extra string, err error) {
	out := fs.Output()

	if isFlagHelpError(err) && fs.Output() == os.Stderr {
		fs.SetOutput(os.Stdout)
	}

	fmt.Fprintf(fs.Output(), "%s:\n", program)
	fs.Usage()
	if extra != "" {
		fmt.Fprintln(fs.Output(), extra)
	}

	fs.SetOutput(out)
}
