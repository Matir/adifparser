package adifparser

import (
	"net/http"
	"strings"
	"testing"
)

type mockBody struct {
	strings.Reader
}

func (_ *mockBody) Close() error {
	return nil
}

func makeMockResponse(s string) *http.Response {
	r := &http.Response{}
	rdr := strings.NewReader(s)
	r.Body = &mockBody{*rdr}
	return r
}

func TestRead(t *testing.T) {
	c := NewLOTWClient("u", "p")
	testString := "This doesn't parse, so nothing special needed.<APP_LoTW_EOF>"
	c.httpResponse = makeMockResponse(testString)
	buf := make([]byte, 1024)
	n, err := c.Read(buf)
	buf = buf[:n]
	if err != nil {
		t.Fatal(err)
	}
	if string(buf[:n]) != testString {
		t.Fatalf("Expected %v, got %v.\n", testString, buf)
	}
}

func TestReadChunked(t *testing.T) {
	c := NewLOTWClient("u", "p")
	testString := "This doesn't parse, so nothing special needed.<APP_LoTW_EOF>"
	c.httpResponse = makeMockResponse(testString)
	buf := make([]byte, 3)
	n, err := c.Read(buf)
	buf = buf[:n]
	if err != nil {
		t.Fatal(err)
	}
	if string(buf[:n]) != testString[:n] {
		t.Fatalf("Expected %v, got %v.\n", testString, buf)
	}
	prev_n := n
	n, err = c.Read(buf)
	buf = buf[:n]
	if err != nil {
		t.Fatal(err)
	}
	if string(buf[:n]) != testString[prev_n:prev_n+n] {
		t.Fatalf("Expected %v, got %v.\n", testString, buf)
	}
}
