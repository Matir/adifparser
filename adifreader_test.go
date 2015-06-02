package adifparser

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func testHeaderFile(t *testing.T, filename string) {
	f, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}

	reader := &baseADIFReader{}
	reader.rdr = f
	reader.readHeader()
	if !bytes.HasPrefix(reader.excess, []byte("<mycall")) {
		t.Fatalf("Excess has %s, expected %s.", string(reader.excess), "<mycall")
	}
}

func TestHeaderNone(t *testing.T) {
	testHeaderFile(t, "testdata/header_none.adi")
}

func TestHeaderVersion(t *testing.T) {
	testHeaderFile(t, "testdata/header_version.adi")
}

func TestHeaderComment(t *testing.T) {
	testHeaderFile(t, "testdata/header_comment.adi")
}

func TestInternalReadRecord(t *testing.T) {
	f, err := os.Open("testdata/readrecord.adi")
	if err != nil {
		t.Fatal(err)
	}

	reader := &baseADIFReader{}
	reader.rdr = f

	testStrings := [...]string{
		"<mycall:6>KF4MDV", "<mycall:6>KG4JEL", "<mycall:4>W1AW"}
	for i := range testStrings {
		buf, err := reader.readRecord()
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}
		if string(buf) != testStrings[i] {
			t.Fatalf("Got bad record %q, expected %q.", string(buf), testStrings[i])
		}
	}
}

func TestReadRecord(t *testing.T) {
	f, err := os.Open("testdata/readrecord.adi")
	if err != nil {
		t.Fatal(err)
	}

	reader := NewADIFReader(f)
	for i := 0; i < 3; i++ {
		_, err = reader.ReadRecord()
		if err != nil && err != io.EOF {
			t.Fatal(err)
		}
	}

	r, err := reader.ReadRecord()
	if err == nil {
		t.Fatal("Expected an error, but err was nil.")
	}
	if err != io.EOF {
		t.Fatalf("Expected %v, got %v", io.EOF, err)
	}
	if r != nil {
		t.Fatalf("Expected nil record, got %v", r)
	}

	_, err = reader.ReadRecord()
	if err == nil {
		t.Fatal("Expected an error, but err was nil.")
	}
	if err != io.EOF {
		t.Fatalf("Expected %v, got %v", io.EOF, err)
	}
}

func TestDedupeReadRecord(t *testing.T) {
	buf := strings.NewReader("<mycall:6>KF4MDV<eor><mycall:6>KF4MDV<eor>")
	reader := NewADIFReader(buf)
	if reader == nil {
		t.Fatal("Invalid reader.")
	}
}
