package adifparser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

var OutputStarted = errors.New("Output already started.")

// Basic writer type
type ADIFWriter interface {
	WriteRecord(ADIFRecord) error
	Flush() error
	SetComment(string) error
}

type baseADIFWriter struct {
	writer  *bufio.Writer
	started bool
}

// Construct a new writer
func NewADIFWriter(w io.Writer) *baseADIFWriter {
	writer := &baseADIFWriter{}
	writer.writer = bufio.NewWriter(w)
	writer.started = false
	return writer
}

func (writer *baseADIFWriter) WriteRecord(r ADIFRecord) error {
	writer.started = true
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

func (writer *baseADIFWriter) SetComment(comment string) error {
	if writer.started {
		return OutputStarted
	}
	fmt.Fprintf(writer.writer, "%s<eoh>\n", comment)
	return nil
}
