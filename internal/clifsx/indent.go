package clifsx

import (
	"bytes"
	"io"
	"reflect"
)

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
