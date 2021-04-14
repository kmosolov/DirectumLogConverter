package DirectumLogConverter

import (
	"bufio"
	jsoniter "github.com/json-iterator/go"
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
		var elements = &MapSlice{}
		var json = jsoniter.ConfigDefault
		_ = json.Unmarshal(raw, &elements)
		logLine := &LogLine{
			Elements:    elements,
			Raw:         raw,
		}
		logEntry := p.converter.Convert(logLine)
		p.printer.Print(logEntry)
	}
	p.printer.Flush()
	return p.scan.Err()
}

type LogLine struct {
	Elements    *MapSlice
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

type LogElement struct {
	Name  string
	Value string
}

type LogEntry struct {
	Elements           []LogElement
	AdditionalElements []string
}
