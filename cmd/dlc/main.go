package main

import (
	"fmt"
	"github.com/kmosolov/DirectumLogConverter"
	flag "github.com/ogier/pflag"
	"os"
	"path/filepath"
	"strings"
)

const VERSION = "0.0.2"
const FilenamePostfix = "_converted"

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

  %s [source] [destination]

  [source] argument is mandatory, [destination] is not, if omitted it will use source file name with postfix "%s" as destination file name.

Switches:

`, filename, filename, FilenamePostfix)
		flag.PrintDefaults()
	}
	var csvFormatArg bool
	var showVersion bool
	flag.BoolVarP(&csvFormatArg, "csv", "c", false, "Use csv as output format.")
	flag.BoolVarP(&showVersion, "version", "v", false, "Print version.")
	flag.Parse()

	if showVersion {
		fmt.Println("DirectumLogConverter. Version: ", VERSION)
		os.Exit(0)
	}

	inFileArg := flag.Arg(0)
	if inFileArg == ""{
		flag.Usage()
		os.Exit(0)
	}

	inFile, err := os.Open(inFileArg)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFileArg := flag.Arg(1)
	if outFileArg == ""{
		fileExt := filepath.Ext(inFileArg)
		outFileArg = fmt.Sprintf("%s%s%s", strings.TrimSuffix(inFileArg, fileExt), FilenamePostfix, fileExt)
	}

	outFile, err := os.Create(outFileArg)
	if err != nil {
		return err
	}
	defer outFile.Close()

	fmt.Fprintf(os.Stdout, "Converting log from \"%s\" to \"%s\"...\n", inFileArg, outFileArg)

	var printer DirectumLogConverter.LogEntryPrinter
	if csvFormatArg  {
		outFile.Write([]byte{0xEF, 0xBB, 0xBF})
		printer = DirectumLogConverter.NewCsvPrinter(outFile)
	} else {
		printer = DirectumLogConverter.NewPrinter(outFile)
	}

	return DirectumLogConverter.NewParser(inFile, printer).Consume()
}
