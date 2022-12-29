package main

import (
	"TestVK/counter"
	"flag"
	"io"
	"log"
	"os"
)

const k = 5

func main() {
	var (
		inputFile  *os.File
		outputFile *os.File
		reader     io.Reader = os.Stdin
		writer     io.Writer = os.Stdout
		fileErr    error
	)

	flag.Parse()

	if inputFileName := flag.Arg(0); inputFileName != "" {
		inputFile, fileErr = os.Open(inputFileName)
		if fileErr != nil {
			log.Fatal(fileErr)
		}
		reader = inputFile
	}

	if outputFileName := flag.Arg(1); outputFileName != "" {
		outputFile, fileErr = os.Create(outputFileName)
		if fileErr != nil {
			log.Fatal(fileErr)
		}
		writer = outputFile
	}

	counter := counter.NewCounter(k, writer)
	err := counter.Execute(reader)
	if err != nil {
		log.Fatal(err)
	}

}
