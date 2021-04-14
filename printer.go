package DirectumLogConverter

import (
	"io"
	"strings"
)

type Printer struct {
	Out io.Writer
}

func NewPrinter(w io.Writer) *Printer {
	return &Printer{Out: w}
}

func (p *Printer) IsWidthFixed() bool {
	return true
}

func (p *Printer) Print(entry *LogEntry) {
	p.Out.Write([]byte(strings.Join(entry.Elements, " ")))
	p.Out.Write([]byte("\n"))
	for _, as := range entry.AdditionalElements {
		p.Out.Write([]byte(as))
		p.Out.Write([]byte("\n"))
	}
}

func (p *Printer) Flush() { }



