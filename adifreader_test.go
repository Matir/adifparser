package adifparser

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testHeaderFile(t *testing.T, filename string) {
	f, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}

	reader := &baseADIFReader{}
	reader.rdr = bufio.NewReader(f)
	reader.readHeader()
	prefix, err := reader.rdr.Peek(128)
	if bytes.HasPrefix(prefix, []byte("<mycall")) {
		t.Fatalf("prefix has %s, expected %s.", string(prefix), "<mycall")
	}
	t.Logf("adif_version: %s", reader.version)
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
	buf := strings.NewReader("<mycall:6>KF4MDV<eor><mycall:6>KF4MDV<fail:1>Y<eor>")
	reader := NewDedupeADIFReader(buf)
	if reader == nil {
		t.Fatal("Invalid reader.")
	}

	if r, err := reader.ReadRecord(); err != nil {
		t.Fatal(err)
	} else if r == nil {
		t.Fatal("Got nil record.")
	}

	if r, err := reader.ReadRecord(); err != io.EOF {
		t.Fatalf("Expected %v, got %v.", io.EOF, err)
	} else if r != nil {
		t.Fatalf("Got %v instead of nil.", r)
	}
}

func TestFullFiles(t *testing.T) {
	expected := map[string]int{
		"lotw.adi":              250,
		"lotw_empty.adi":        0,
		"lotw_empty_no_eof.adi": 0,
		"lotw_eof.adi":          1,
		"lotw_new.adi":          4,
		"wsjtx.adi":             74,
		"xlog.adi":              425,
	}
	for file, count := range expected {
		fname := filepath.Join("testdata", file)
		f, err := os.Open(fname)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		reader := NewADIFReader(f)
		for rec, err := reader.ReadRecord(); err != io.EOF; {
			if err != nil {
				t.Fatalf("input %s: %v", file, err)
			}
			if rec == nil {
				t.Fatalf("input %s: rec is nil", file)
			}
			rec, err = reader.ReadRecord()
		}
		if reader.RecordCount() != count {
			t.Fatalf("Record count for %s was wrong: got %d, expected %d.",
				file, reader.RecordCount(), count)
		}
	}
}

func TestReadRecordWithBlank(t *testing.T) {
	buf := strings.NewReader("<mycall:7>KF4MDV <eor>")
	reader := NewADIFReader(buf)
	if reader == nil {
		t.Fatal("Invalid reader.")
	}

	r, err := reader.ReadRecord()
	if err != nil {
		t.Fatal(err)
	}
	if r == nil {
		t.Fatal("Got nil record.")
	}

	if v, err := r.GetValue("mycall"); err != nil {
		t.Fatal("Got value:mycall error")
	} else {
		if v != "KF4MDV " {
			t.Fatal("Not matched")
		}
	}

}

func TestReadRecordWithBlank2(t *testing.T) {
	buf := strings.NewReader(" <EOH>\n<NOTES:10>          <eor>")
	reader := NewADIFReader(buf)
	if reader == nil {
		t.Fatal("Invalid reader.")
	}

	r, err := reader.ReadRecord()
	if err != nil {
		t.Fatal(err)
	}
	if r == nil {
		t.Fatal("Got nil record.")
	}

	if v, err := r.GetValue("notes"); err != nil {
		t.Fatal("Got value:notes error")
	} else {
		if v != "          " {
			t.Fatal("Not matched")
		}
	}

}

func TestReadRecordWithFiller(t *testing.T) {
	buf := strings.NewReader("<NOTES:8>ABCDEFGH|IGNORE_THIS|<EOR>|FILLER")
	reader := NewADIFReader(buf)
	if reader == nil {
		t.Fatal("Invalid reader.")
	}

	r, err := reader.ReadRecord()
	if err != nil {
		t.Fatal(err)
	}
	if r == nil {
		t.Fatal("Got nil record.")
	}

	if v, err := r.GetValue("notes"); err != nil {
		t.Fatal("Got value:notes error")
	} else {
		if v != "ABCDEFGH" {
			t.Fatal("Not matched")
		}
	}

}

// This is out of ADI file specification,
// but in some cases you should expect non-ASCII characters in the field value.

func TestReadRecordWithNonASCII(t *testing.T) {
	buf := strings.NewReader("<TEXT:4>AB\xedD\xeb<EOR> ")
	reader := NewADIFReader(buf)
	if reader == nil {
		t.Fatal("Invalid reader.")
	}

	r, err := reader.ReadRecord()
	if err != nil {
		t.Fatal(err)
	}
	if r == nil {
		t.Fatal("Got nil record.")
	}

	if v, err := r.GetValue("text"); err != nil {
		t.Fatal("Got value:text error")
	} else {
		if v != "AB\xedD" {
			t.Fatal("Not matched")
		}
	}

}

func TestReadRecordWithNonASCII2(t *testing.T) {
	buf := strings.NewReader("<noTes:5>12345<test:6> XY\xedZ <yum:7:s>ABCD!EF<EOR>")
	reader := NewADIFReader(buf)
	if reader == nil {
		t.Fatal("Invalid reader.")
	}

	r, err := reader.ReadRecord()
	if err != nil {
		t.Fatal(err)
	}
	if r == nil {
		t.Fatal("Got nil record.")
	}

	if v, err := r.GetValue("notes"); err != nil {
		t.Fatal("Got value:notes error")
	} else {
		if v != "12345" {
			t.Fatal("Not matched")
		}
	}
	if v, err := r.GetValue("test"); err != nil {
		t.Fatal("Got value:test error")
	} else {
		if v != " XY\xedZ " {
			t.Fatal("Not matched")
		}
	}
	if v, err := r.GetValue("yum"); err != nil {
		t.Fatal("Got value:yum error")
	} else {
		if v != "ABCD!EF" {
			t.Fatal("Not matched")
		}
	}

}

func TestReadRecordWithNoEOH(t *testing.T) {
	buf := strings.NewReader(" <TEST:1>A <EOR> ")
	reader := NewADIFReader(buf)
	if reader == nil {
		t.Fatal("Invalid reader.")
	}

	_, err := reader.ReadRecord()
	if err != io.EOF {
		t.Fatalf("Expected %v, got %v", io.EOF, err)
	}
}

func TestForReadElement(t *testing.T) {
	buf := strings.NewReader(" |FILLER1| <TeSt:2>XY |FILLER 2| <eOr>  ")
	reader := NewADIFReader(buf)
	if reader == nil {
		t.Fatal("Invalid reader.")
	}
	element, err := reader.readElement()
	if err != nil {
		t.Fatal(err)
	}
	if element == nil {
		t.Fatal("Got nil element.")
	}
	if element.name != "test" {
		t.Fatal("element.name not matched for TEST")
	}
	if element.value != "XY" {
		t.Fatal("element.value not matched for XY")
	}
	element, err = reader.readElement()
	if err != nil {
		t.Fatal(err)
	}
	if element == nil {
		t.Fatal("Got nil element.")
	}
	if element.name != "eor" {
		t.Fatal("element.name not matched")
	}
	if element.hasValue {
		t.Fatal("element.value for EOR is true")
	}

}
