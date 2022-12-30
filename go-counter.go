package main

import (
	"flag"
	"fmt"
	"github.com/z-urvancev/go-counter/counter"
	"io"
	"log"
	"os"
)

const k = 5

func main() {
	var (
		inputFile      *os.File
		outputFile     *os.File
		reader         io.Reader = os.Stdin
		writer         io.Writer = os.Stdout
		fileErr        error
		helpFlag       bool
		inputFileName  string
		outputFileName string
	)

	flag.StringVar(&inputFileName, "input", "", "input file. By default use stdin")
	flag.StringVar(&outputFileName, "output", "", "output file. By default use stdout")
	flag.BoolVar(&helpFlag, "help", false, "help for go-counter")
	flag.Parse()

	if helpFlag {
		fmt.Println("Usage:\ngo-counter [--help] [--input=<input_file_name>] [--output=<output_file_name>]\nFlags:\n--help - help for go-counter\n--input string - input file. By default use stdin\n--output string - output file. By default use stdout")
		return
	}

	if inputFileName != "" {
		inputFile, fileErr = os.Open(inputFileName)
		if fileErr != nil {
			log.Fatalln("cannot open input file: ", fileErr)
		}
		reader = inputFile
	}

	if outputFileName != "" {
		outputFile, fileErr = os.Create(outputFileName)
		if fileErr != nil {
			log.Fatalln("cannot open output file: ", fileErr)
		}
		writer = outputFile
	}

	counter := counter.NewCounter(k, writer)
	err := counter.Execute(reader)
	if err != nil {
		log.Fatal(err)
	}

}
