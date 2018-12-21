# clip

    go get -u github.com/codemodus/clip/clipr


Package clipr provides error implementations and helpers related to the clip and 
clifs packages.

## Usage

```go
func FilterControl(err error) error
func IsFlagHelp(err error) bool
type BadCommandError
    func NewBadCommandError(scope, cmd string) *BadCommandError
    func (e *BadCommandError) Error() string
type EmptyCommandError
    func NewEmptyCommandError(scope string) *EmptyCommandError
    func (e *EmptyCommandError) Error() string
type FlagParseError
    func NewFlagParseError(scope string, err error) *FlagParseError
    func (e *FlagParseError) Error() string
type UsageError
    func NewUsageError(err error, usageFunc func(int, error) error, top bool) *UsageError
    func (e *UsageError) Err() error
    func (e *UsageError) Error() string
    func (e *UsageError) IsHelp() bool
    func (e *UsageError) IsTop() bool
    func (e *UsageError) Usage(depth int) error
```

### Setup

N/A

## More Info

The clip package is setup to provide a reasonably large amount of flexibility.
This package is being made available to accommodate additional potentialities. 

## Documentation

View the [GoDoc](http://godoc.org/github.com/codemodus/clip/clipr)

## Benchmarks

N/A
