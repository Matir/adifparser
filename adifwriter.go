package adifparser

import (
	"bufio"
	"fmt"
	"io"
)

// Basic writer type
type ADIFWriter interface {
	WriteRecord(ADIFRecord) error
	Flush() error
}

type baseADIFWriter struct {
	writer *bufio.Writer
}

// Construct a new writer
func NewADIFWriter(w io.Writer) *baseADIFWriter {
	writer := &baseADIFWriter{}
	writer.writer = bufio.NewWriter(w)
	return writer
}

func (writer *baseADIFWriter) WriteRecord(r ADIFRecord) error {
	_, err := fmt.Fprintf(writer.writer, "%s<eor>\n", r.ToString())
	if err != nil {
		// TODO: log
		return err
	}
	return nil
}

func (writer *baseADIFWriter) Flush() error {
	return writer.writer.Flush()
}
