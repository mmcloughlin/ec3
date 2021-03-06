package print

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

const DefaultIndent = "\t"

type Printer struct {
	out     io.Writer
	level   int    // current indentation level
	indent  string // indentation string
	pending bool   // if there's a pending indentation
	err     error  // saved error from printing
}

func New(w io.Writer) Printer {
	return Printer{
		out:    w,
		indent: DefaultIndent,
	}
}

func (p *Printer) SetIndentString(indent string) {
	p.indent = indent
}

func (p *Printer) Indent() {
	p.level++
}

func (p *Printer) Dedent() {
	p.level--
}

func (p *Printer) Linef(format string, args ...interface{}) {
	p.Printf(format, args...)
	p.NL()
}

func (p *Printer) NL() {
	p.Printf("\n")
	p.pending = true
}

func (p *Printer) Printf(format string, args ...interface{}) {
	if p.err != nil {
		return
	}
	if p.pending {
		indent := strings.Repeat(p.indent, p.level)
		format = indent + format
		p.pending = false
	}
	_, err := fmt.Fprintf(p.out, format, args...)
	p.SetError(err)
}

func (p *Printer) Error() error {
	return p.err
}

func (p *Printer) SetError(err error) {
	if p.err == nil {
		p.err = err
	}
}

// Buffer is a convenience wrapper for printing into a byte buffer.
type Buffer struct {
	Buf *bytes.Buffer
	Printer
}

func NewBuffer() Buffer {
	buf := bytes.NewBuffer(nil)
	return Buffer{
		Buf:     buf,
		Printer: New(buf),
	}
}

func (b *Buffer) Result() ([]byte, error) {
	return b.Buf.Bytes(), b.Error()
}

type TabWriter struct {
	tw *tabwriter.Writer
	Printer
}

func NewTabWriter(w io.Writer, minwidth, tabwidth, padding int, padchar byte, flags uint) *TabWriter {
	tw := tabwriter.NewWriter(w, minwidth, tabwidth, padding, padchar, flags)
	return &TabWriter{
		tw:      tw,
		Printer: New(tw),
	}
}

func (p *TabWriter) Flush() {
	p.SetError(p.tw.Flush())
}
