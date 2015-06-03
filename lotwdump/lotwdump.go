package main

import (
	"flag"
	"fmt"
	"github.com/Matir/adifparser"
	"io"
	"os"
)

func main() {
	var username = flag.String("username", "", "LOTW Username")
	var password = flag.String("password", "", "LOTW Password")

	flag.Parse()

	if *username == "" || *password == "" {
		fmt.Fprintf(os.Stderr, "Username and password required.\n")
		return
	}

	client := adifparser.NewLOTWClient(*username, *password)
	reader := adifparser.NewADIFReader(client)
	writer := adifparser.NewADIFWriter(os.Stdout)
	defer writer.Flush()

	for record, err := reader.ReadRecord(); record != nil; record, err = reader.ReadRecord() {
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}
			return
		}
		writer.WriteRecord(record)
	}
}
