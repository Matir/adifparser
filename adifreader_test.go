package adifparser

import (
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
