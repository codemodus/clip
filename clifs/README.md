# clifs

    go get -u github.com/codemodus/clip/clifs

Package clifs wraps flag.FlagSet to improve the flexibility of usage output. If 
the FlagSet output is default, setting a help flag will result in usage output 
going to os.Stdout. Other parse errors will result in usage output going to 
os.Stderr as normal.

## Usage

```go
func Parse(fs *FlagSet, args []string) error
func Usage(program string, fs *FlagSet, extra string, err error) error
type FlagSet
    func NewFlagSet(name string) *FlagSet
```

### Setup

```go
var (
    program = path.Base(os.Args[0])
    verbose bool
    example string
)

fs := clifs.NewFlagSet("global")
fs.BoolVar(&verbose, "v", verbose, "verbosity")
fs.StringVar(&example, "example", example, "example")

if err := clifs.Parse(fs, os.Args[1:]); err != nil {
    return clifs.Usage(program, fs, "", err)
}

fmt.Println(verbose, example)
return nil
```

## More Info

N/A

## Documentation

View the [GoDoc](http://godoc.org/github.com/codemodus/clip/clifs)

## Benchmarks

N/A
