package main

import (
	"flag"
	"fmt"
	"github.com/Matir/adifparser"
	"io"
	"os"
	"time"
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
	defer client.Close()

	t := time.Now().Format("2006/01/02 15:04:05")
	writer.SetComment(fmt.Sprintf("Downloaded from LOTW at %s.", t))

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
