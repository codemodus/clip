# clip

    go get -u github.com/codemodus/clip

Package clip provides simple command/subcommand structuring with a minimum of 
dependencies. The standard library flag.FlagSet type is leveraged often in this 
and related packages. If the FlagSet output is default, setting a help flag will 
result in usage output going to os.Stdout. Other parse errors will result in 
usage output going to os.Stderr as normal.

## Usage

```go
type Clip
    func New(program string, flags *FlagSet, subcmds *CommandSet) *Clip
    func (c *Clip) Parse(args []string) error
    func (c *Clip) Run() error
    func (c *Clip) Usage(depth int, err error) error
    func (c *Clip) UsageLongHelp(err error) error
type Command
    func NewCommand(flags *FlagSet, fn HandlerFunc, subcmds *CommandSet) *Command
    func NewCommandNamespace(name string, subcmds *CommandSet) *Command
type CommandSet
    func NewCommandSet(cmds ...*Command) *CommandSet
type FlagSet
    func NewFlagSet(name string) *FlagSet
type HandlerFunc
type UsageError
    func AsUsageError(err error) (UsageError, bool)
```

### Setup

```go
var (
    globalCnf = newGlobalConf()
    printCnf  = newPrintConf("print")
    otherCnf  = newOtherConf("other")
)

cs := clip.NewCommandSet(
    clip.NewCommand(printCnf.flagSet, runPrintFunc(printCnf, globalCnf), nil),
    clip.NewCommand(otherCnf.flagSet, runOtherFunc(otherCnf), nil),
)
app := clip.New("myapp", globalCnf.flagSet, cs)

// emulate cli command 'myapp -v print -msg=hello, world'
os.Args = []string{"myapp", "-v", "print", "-msg=hello, world"}

if err := app.Parse(os.Args); err != nil {
    return app.UsageLongHelp(err)
}

return app.Run()
```

## More Info

N/A

## Documentation

View the [GoDoc](http://godoc.org/github.com/codemodus/clip)

## Benchmarks

N/A
