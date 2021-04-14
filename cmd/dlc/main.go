package main

import (
	"fmt"
	"github.com/kmosolov/DirectumLogConverter"
	flag "github.com/ogier/pflag"
	"os"
	"path/filepath"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	flag.Usage = func() {
		filename := filepath.Base(os.Args[0])
		fmt.Printf(`DirectumLogConverter is a tool for converting Directum JSON logs.

Usage of %s:

  %s [filename]

  If [filename] is omitted, it reads from standard input.

Switches:

`, filename, filename)
		flag.PrintDefaults()
	}
	var csvFormatArg bool
	var outFileArg string
	flag.BoolVarP(&csvFormatArg, "csv", "c", false, "Use csv as output format.")
	flag.StringVarP(&outFileArg, "output", "o", "", "Output file, if omitted it writes to standard output.")
	flag.Parse()

	inFileArg := flag.Arg(0)
	inFile := os.Stdin
	if inFileArg != "" {
		f, err := os.Open(inFileArg)
		if err != nil {
			return err
		}
		defer f.Close()
		inFile = f
	}

	outFile := os.Stdout
	if outFileArg != "" {
		f, err := os.Create(outFileArg)
		if err != nil {
			return err
		}
		defer f.Close()
		outFile = f
	}

	if outFileArg != "" {
		var inputStr string
		if inFileArg == "" {
			inputStr = "standard input"
		} else {
			inputStr = fmt.Sprintf("\"%s\"", inFileArg)
		}
		fmt.Fprintf(os.Stdout, "Converting log from %s to \"%s\"...\n", inputStr, outFileArg)
	}

	var printer DirectumLogConverter.LogEntryPrinter
	if csvFormatArg  {
		printer = DirectumLogConverter.NewCsvPrinter(outFile)
	} else {
		printer = DirectumLogConverter.NewPrinter(outFile)
	}

	return DirectumLogConverter.NewParser(inFile, printer).Consume()
}
