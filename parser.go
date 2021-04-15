package DirectumLogConverter

import (
	"bufio"
	jsoniter "github.com/json-iterator/go"
	"io"
	"runtime"
	"sync"
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

func processLinesArray(p *Parser, a *[]*LogLine, len int) {
	var wg sync.WaitGroup
	wg.Add(len)
	for i := 0; i < len; i++ {
		go processLine(p, (*a)[i], &wg)
	}
	wg.Wait()
	for i := 0; i < len; i++ {
		p.printer.Print((*a)[i].Entry)
	}
}

func processLine(p *Parser, logLine *LogLine, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	logLine.Elements = &MapSlice{}
	var json = jsoniter.ConfigDefault
	_ = json.Unmarshal(*logLine.Raw, logLine.Elements)
	logLine.Entry = p.converter.Convert(logLine)
}

func parallelProcessLines(p *Parser) {
	buf := make([]*LogLine, runtime.NumCPU()*2)
	i := 0
	scan := p.scan
	for scan.Scan() {
		bytes := scan.Bytes()
		raw := make([]byte, len(bytes))
		copy(raw, bytes)
		if i == len(buf) {
			processLinesArray(p, &buf, i)
			i = 0
		}
		buf[i] = &LogLine{Raw: &raw}
		i++
	}
	if i > 0 {
		processLinesArray(p, &buf, i)
	}
}

func serialProcessLines(p *Parser) {
	scan := p.scan
	for scan.Scan() {
		raw := scan.Bytes()
		logLine := &LogLine{Raw: &raw}
		processLine(p, logLine, nil)
		p.printer.Print(logLine.Entry)
	}
}

func (p *Parser) Consume(canParallelize bool) error {
	if canParallelize {
		parallelProcessLines(p)
	} else {
		serialProcessLines(p)
	}
	p.printer.Flush()
	return p.scan.Err()
}

type LogLine struct {
	Elements *MapSlice
	Raw      *[]byte
	Entry    *LogEntry
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
	Elements           *[]*LogElement
	AdditionalElements *[]string
}
