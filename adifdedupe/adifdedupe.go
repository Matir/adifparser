package main

import (
	"flag"
	"fmt"
	"github.com/Matir/adifparser"
	"io"
	"os"
)

func main() {
	var infile = flag.String("infile", "", "Input file.")
	var outfile = flag.String("outfile", "", "Output file.")

	flag.Parse()

	if *infile == "" {
		fmt.Fprint(os.Stderr, "Need infile.\n")
		return
	}

	fp, err := os.Open(*infile)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}

	var writer adifparser.ADIFWriter
	var writefp *os.File
	if *outfile != "" {
		writefp, err = os.Create(*outfile)
		writer = adifparser.NewADIFWriter(writefp)
	} else {
		writefp = nil
		writer = adifparser.NewADIFWriter(os.Stdout)
	}

	reader := adifparser.NewDedupeADIFReader(fp)
	for record, err := reader.ReadRecord(); record != nil || err != nil; record, err = reader.ReadRecord() {
		if err != nil {
			if err != io.EOF {
				fmt.Fprint(os.Stderr, err)
			}
			break
		}
		writer.WriteRecord(record)
	}

	if writefp != nil {
		writefp.Close()
	}
}
