package adifparser

import (
	"testing"
)

func TestGetNextField(t *testing.T) {
	buf := []byte("<blah:2>AB<FOO:3>XYZ <bar:4:s>1234")

	expected := []struct {
		n string
		v string
	}{
		{"blah", "AB"},
		{"foo", "XYZ"},
		{"bar", "1234"},
	}

	var err error
	var data *fieldData

	for _, el := range expected {
		data, buf, err = getNextField(buf)
		if err != nil {
			t.Fatal(err)
		}
		if data.name != el.n || data.value != el.v {
			t.Fatalf("Got %q=%q, expected %q=%q.", data.name, data.value, el.n, el.v)
		}
	}
}
