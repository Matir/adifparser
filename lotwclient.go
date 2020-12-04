package adifparser

import (
	"bytes"
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
	// Options
	qsl_only      bool
	qso_mydetail  bool
	qso_qsldetail bool
	qso_withown   bool
	// TODO: state storage
}

// Create a new client
func NewLOTWClient(username, password string) *lotwClientImpl {
	client := &lotwClientImpl{}
	client.username = username
	client.password = password
	client.buf = make([]byte, 0, 1024)

	// Default options
	client.qsl_only = false
	client.qso_mydetail = true
	client.qso_qsldetail = true
	client.qso_withown = true

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
		n = copy(p, usable)
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

	n = copy(p, usable[:max_len])
	if err == io.EOF && len(c.buf) == 0 {
		return n, err
	}
	return n, nil
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
	params := c.getParams()
	requri, _ := url.Parse(LOTWAPI)
	requri.RawQuery = makeQueryString(params)
	adiflog.Printf("LOTW Requesting: %s\n", requri.String())
	resp, err := http.Get(requri.String())
	if err != nil {
		return err
	}
	c.httpResponse = resp
	return nil
}

func (c *lotwClientImpl) getParams() map[string]string {
	params := map[string]string{
		"login":     c.username,
		"password":  c.password,
		"qso_query": "1",
	}
	if c.qsl_only {
		params["qso_qsl"] = "yes"
	} else {
		params["qso_qsl"] = "no"
	}
	if c.qso_mydetail {
		params["qso_mydetail"] = "yes"
	}
	if c.qso_qsldetail {
		params["qso_qsldetail"] = "yes"
	}
	if c.qso_withown {
		params["qso_withown"] = "yes"
	}
	return params
}

func makeQueryString(data map[string]string) string {
	elements := make([]string, 0, len(data))
	for param, value := range data {
		elements = append(elements,
			url.QueryEscape(param)+"="+url.QueryEscape(value))
	}
	return strings.Join(elements, "&")
}
