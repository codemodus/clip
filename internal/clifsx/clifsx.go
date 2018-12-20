package clifsx

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"reflect"

	"github.com/codemodus/clip/clipr"
)

// Parse ...
func Parse(fs *flag.FlagSet, args []string) error {
	args = shiftCollision(args, fs.Name(), os.Args[0], path.Base(os.Args[0]))

	out := fs.Output()
	fs.SetOutput(ioutil.Discard)
	defer fs.SetOutput(out)

	return fs.Parse(args)
}

// Usage ...
func Usage(program string, depth int, fs *flag.FlagSet, extra string, err error) error {
	if err == nil {
		return nil
	}

	if clipr.IsFlagHelp(err) {
		err = nil

		if fs.Output() == os.Stderr {
			fs.SetOutput(os.Stdout)
			defer fs.SetOutput(nil)
		}
	}

	if depth > 0 {
		it := newIndentTool(fs.Output(), "    ", depth)
		fs.SetOutput(it)
		defer fs.SetOutput(it.w)

		fmt.Fprintln(fs.Output())
	}

	if program != "" && program != fs.Name() {
		fmt.Fprintf(fs.Output(), "%s:\n", program)
	}

	fs.Usage()

	if extra != "" {
		fmt.Fprintln(fs.Output(), extra)
	}

	return err
}

func shiftCollision(args []string, ss ...string) []string {
	for _, s := range ss {
		if len(args) > 0 && args[0] == s {
			args = args[1:]
			break
		}
	}

	return args
}

type indentTool struct {
	ind []byte
	trg []byte
	alt []byte
	w   io.Writer
}

func newIndentTool(w io.Writer, indent string, depth int) *indentTool {
	ind := bytes.Repeat([]byte(indent), depth)
	trg := []byte("\n")
	alt := append(trg, ind...)

	return &indentTool{
		ind: ind,
		trg: trg,
		alt: alt,
		w:   w,
	}
}

func (i *indentTool) Write(p []byte) (n int, err error) {
	bs := i.ind
	rp := bytes.Replace(p, i.trg, i.alt, -1)

	bs = append(bs, rp...)

	if reflect.DeepEqual(bs[len(bs)-len(i.alt):], i.alt) {
		bs = bs[:len(bs)-len(i.alt)+len(i.trg)]
	}

	return i.w.Write(bs)
}
