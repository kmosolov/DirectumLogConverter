package DirectumLogConverter

import (
	"encoding/csv"
	"io"
	"log"
	"strings"
)

var mandatoryLogElements = []string{"t", "pid", "v", "un", "tn"}

type CsvPrinter struct {
	Writer *csv.Writer
}

func NewCsvPrinter(w io.Writer) *CsvPrinter {
	writer := csv.NewWriter(w)
	writer.Comma = ';'
	writer.UseCRLF = true
	return &CsvPrinter{Writer: writer}
}

func (p *CsvPrinter) IsWidthFixed() bool {
	return false
}

func (p *CsvPrinter) Print(entry *LogEntry) {
	var visitedElements = make(map[string]bool)
	var record []string
	for _, elementName := range mandatoryLogElements {
		value := ""
		for _, element := range entry.Elements {
			if element.Name == elementName {
				value = element.Value
				break
			}
		}
		visitedElements[elementName] = true
		record = append(record, value)
	}
	for _, element := range entry.Elements {
		if !visitedElements[element.Name] {
			record = append(record, element.Value)
		}
	}
	if len(entry.AdditionalElements) > 0 {
		record = append(record, strings.Join(entry.AdditionalElements, "\n"))
	}
	if err := p.Writer.Write(record); err != nil {
		log.Fatal(err)
	}
}

func (p *CsvPrinter) Flush() {
	p.Writer.Flush()

	if err := p.Writer.Error(); err != nil {
		log.Fatal(err)
	}
}
