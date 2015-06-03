package adifparser

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const LOTWAPI string = "https://lotw.arrl.org/lotwuser/lotwreport.adi"

type LOTWClient interface {
	// Implement reader
	Read([]byte) (int, error)
	Close() error
}

type lotwClientImpl struct {
	// Creds
	username string
	password string
	// HTTP state
	httpResponse *http.Response
	// Temporary read buffer
	buf []byte
	// TODO: state storage
}

// Create a new client
func NewLOTWClient(username, password string) *lotwClientImpl {
	client := &lotwClientImpl{}
	client.username = username
	client.password = password
	client.buf = make([]byte, 0, 1024)
	return client
}

// Read from socket
func (c *lotwClientImpl) Read(p []byte) (int, error) {
	if c.httpResponse == nil {
		if err := c.open(); err != nil {
			// TODO: better logging
			return 0, err
		}
	}
	usable := c.buf
	storage := make([]byte, 1024)
	n, err := c.httpResponse.Body.Read(storage)
	if err != nil && err != io.EOF {
		n = copyWithoutLOTW(p, usable)
		return n, err
	}
	usable = append(usable, storage[:n]...)
	max_len := len(usable)
	if cap(p) < max_len {
		max_len = cap(p)
	}
	if end := bytes.LastIndex(usable[:max_len], []byte("<")); end > 0 {
		if endtag := bytes.LastIndex(usable, []byte(">")); endtag > -1 && endtag < end {
			max_len = end
		}
	}

	c.buf = usable[max_len:]

	n = copyWithoutLOTW(p, usable[:max_len])
	if err == io.EOF && len(c.buf) == 0 {
		return n, err
	}
	return n, nil
}

func copyWithoutLOTW(dst, src []byte) int {
	if end := bytes.Index(src, []byte("<APP_LoTW_EOF>")); end > -1 {
		src = src[:end]
	}
	return copy(dst, src)
}

func (c *lotwClientImpl) Close() error {
	if c.httpResponse != nil {
		err := c.httpResponse.Body.Close()
		c.httpResponse = nil
		return err
	}
	return nil
}

func (c *lotwClientImpl) open() error {
	params := map[string]string{
		"login":     c.username,
		"password":  c.password,
		"qso_query": "1",
	}
	requri, _ := url.Parse(LOTWAPI)
	requri.RawQuery = makeQueryString(params)
	fmt.Printf("Requesting: %s\n", requri.String())
	resp, err := http.Get(requri.String())
	if err != nil {
		return err
	}
	c.httpResponse = resp
	return nil
}

func makeQueryString(data map[string]string) string {
	elements := make([]string, 0, len(data))
	for param, value := range data {
		elements = append(elements,
			url.QueryEscape(param)+"="+url.QueryEscape(value))
	}
	return strings.Join(elements, "&")
}
