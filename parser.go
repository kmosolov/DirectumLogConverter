package DirectumLogConverter

import (
	"bufio"
	"encoding/json"
	"github.com/mickep76/mapslice-json"
	"io"
)

type Parser struct {
	r         io.Reader
	scan      *bufio.Scanner
	converter LogLineConverter
	printer   LogEntryPrinter
}

func NewParser(r io.Reader, printer LogEntryPrinter) *Parser {
	return &Parser{
		r:         r,
		scan:      bufio.NewScanner(r),
		converter: NewConverter(printer.IsWidthFixed()),
		printer:   printer,
	}
}

func (p *Parser) Consume() error {
	s := p.scan
	for s.Scan() {
		raw := s.Bytes()
		var elements = mapslice.MapSlice{}
		_ = json.Unmarshal(raw, &elements)
		logLine := &LogLine{
			Elements:    elements,
			Raw:         raw,
		}
		logEntry := p.converter.Convert(logLine)
		p.printer.Print(logEntry)
	}
	return p.scan.Err()
}

type LogLine struct {
	Elements    mapslice.MapSlice
	Raw         []byte
}

type LogLineConverter interface {
	Convert(*LogLine) *LogEntry
}

type LogEntryPrinter interface {
	IsWidthFixed() bool
	Print(*LogEntry)
	Flush()
}

type LogEntry struct {
	Elements           []string
	AdditionalElements []string
}
