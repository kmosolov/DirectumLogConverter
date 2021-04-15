package main

import (
	"bufio"
	"fmt"
	"github.com/kmosolov/DirectumLogConverter"
	flag "github.com/ogier/pflag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const VERSION = "0.0.3"
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
	var pipeArg bool
	var csvFormatArg bool
	var showVersion bool
	flag.BoolVarP(&pipeArg, "pipe", "p", false, "Pipeline mode, input from STDIN, output to STDOUT.")
	flag.BoolVarP(&csvFormatArg, "csv", "c", false, "Use csv as output format.")
	flag.BoolVarP(&showVersion, "version", "v", false, "Print version.")
	flag.Parse()

	if showVersion {
		fmt.Println("DirectumLogConverter. Version: ", VERSION)
		os.Exit(0)
	}

	inFileArg := flag.Arg(0)
	outFileArg := flag.Arg(1)
	if inFileArg == "" && !pipeArg {
		flag.Usage()
		os.Exit(0)
	}

	var input io.Reader
	var output io.Writer

	if pipeArg {
		if inFileArg != "" || outFileArg != "" {
			fmt.Fprintln(os.Stderr, "[source] and [destination] arguments cannot be used in pipeline mode.")
			os.Exit(1)
		}
		input = os.Stdin
		output = os.Stdout
	} else {
		inFile, err := os.Open(inFileArg)
		if err != nil {
			return err
		}
		input = inFile
		defer inFile.Close()

		if outFileArg == "" {
			fileExt := filepath.Ext(inFileArg)
			newFileExt := fileExt
			if csvFormatArg {
				newFileExt = ".csv"
			}
			outFileArg = fmt.Sprintf("%s%s%s", strings.TrimSuffix(inFileArg, fileExt), FilenamePostfix, newFileExt)
		}

		if fileExists(outFileArg) && !askForConfirmation(fmt.Sprintf("File \"%s\" already exists, overwrite?", outFileArg)) {
			os.Exit(0)
		}

		outFile, err := os.Create(outFileArg)
		if err != nil {
			return err
		}
		output = outFile
		defer outFile.Close()

		fmt.Fprintf(os.Stdout, "Converting log from \"%s\" to \"%s\"...\n", inFileArg, outFileArg)
	}

	var printer DirectumLogConverter.LogEntryPrinter
	if csvFormatArg {
		if !pipeArg {
			output.Write([]byte{0xEF, 0xBB, 0xBF})
		}
		printer = DirectumLogConverter.NewCsvPrinter(output)
	} else {
		printer = DirectumLogConverter.NewPrinter(output)
	}

	start := time.Now()
	result := DirectumLogConverter.NewParser(input, printer).Consume(!pipeArg)
	if !pipeArg {
		fmt.Fprintf(os.Stdout, "Done! %s elapsed.", time.Since(start))
	}
	return result
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}

		fmt.Printf("Unrecognized input \"%s\"\n", response)
	}
}
