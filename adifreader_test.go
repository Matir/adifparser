package adifparser

import (
	"bytes"
	"os"
	"testing"
)

func testHeaderFile(t *testing.T, filename string) {
	f, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}

	reader := new(baseAdifReader)
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
