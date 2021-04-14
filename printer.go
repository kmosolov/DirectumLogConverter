package DirectumLogConverter

import (
	"bufio"
	"io"
)

type Printer struct {
	Writer *bufio.Writer
}

func NewPrinter(w io.Writer) *Printer {
	return &Printer{bufio.NewWriter(w)}
}

func (p *Printer) IsWidthFixed() bool {
	return true
}

func (p *Printer) Print(entry *LogEntry) {
	firstElement := true
	for _, element := range entry.Elements {
		if firstElement {
			firstElement = false
		} else {
			p.Writer.WriteRune(' ')
		}
		p.Writer.WriteString(element.Value)
	}

	p.Writer.WriteRune('\n')
	for _, additionalElement := range entry.AdditionalElements {
		p.Writer.WriteString(additionalElement)
		p.Writer.WriteRune('\n')
	}
}

func (p *Printer) Flush() {
	p.Writer.Flush()
}



