package DirectumLogConverter

import (
	"encoding/csv"
	"io"
	"log"
)

type CsvPrinter struct {
	Writer *csv.Writer
}

func NewCsvPrinter(w io.Writer) *CsvPrinter {
	writer := csv.NewWriter(w)
	writer.Comma = ';'
	return &CsvPrinter{Writer: writer}
}

func (p *CsvPrinter) IsWidthFixed() bool {
	return false
}

func (p *CsvPrinter) Print(entry *LogEntry) {
	if err := p.Writer.Write(append(entry.Elements, entry.AdditionalElements...)); err != nil {
		log.Fatal(err)
	}
}

func (p *CsvPrinter) Flush() {
	p.Writer.Flush()

	if err := p.Writer.Error(); err != nil {
		log.Fatal(err)
	}
}
